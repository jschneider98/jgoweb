package jgoweb

import (
	"fmt"
	"time"
	"os"
	"context"
	"net/http"
	"crypto/tls"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gocraft/web"
	"github.com/gocraft/health"
	"github.com/alexedwards/scs"
	"github.com/jschneider98/jgoweb/config"
)

var healthStream = health.NewStream()
var sessionManager *scs.Manager
var appConfig *config.Config
var appConfigPath string = "./config/config.json"

//
func SetConfigPath(path string) {
	appConfigPath = path
}

//
func GetConfigPath() string {
	return appConfigPath
}

// Init config
func InitConfig() {
	var err error

	if appConfig != nil {
		return
	}

	appConfig, err = config.New(appConfigPath)

	if err != nil {
		panic(err)
	}
}

//
func GetAppConfig() *config.Config {
	return appConfig
}

// Init session
func InitSession() {
	InitConfig()

	sessionManager = scs.NewCookieManager(appConfig.Server.SessionKey)
	scs.CookieName = appConfig.Server.SessionName
}

//
func Start(router *web.Router) {
	InitConfig()
	StartAll(router)
}

//
func StartAll(router *web.Router) {
	InitConfig()
	InitDbCollection()
	InitSession()
	StartHealthSink(appConfig.Server.HealthHost)

	if appConfig.Server.EnableSsl {
		StartHttpsServer(router)
	} else {
		StartHttpServer(router, appConfig.Server.HttpHost)
	}
}

//
func GetWebServer(router *web.Router, host string) *http.Server  {

	server := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      router,
		}

	server.Addr = host

	return server
}

// Start Health Sink
func StartHealthSink(hostname string) {
	healthStream.AddSink(&health.WriterSink{os.Stdout})
	sink := health.NewJsonPollingSink(time.Minute, time.Minute*5)
	healthStream.AddSink(sink)

	fmt.Println("Health Sink Running: ", hostname)
	sink.StartServer(hostname)
}

//
func StartHttpServer(router *web.Router, host string) {
	server := GetWebServer(router, host)

	fmt.Println("HTTP Server Running: ", server.Addr)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

// Start a HTTPS server that auto updates SSL certs via ACME
func StartHttpsServer(router *web.Router) {
	InitConfig()

	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := appConfig.Autocert.AllowedHost

		if host == allowedHost {
			return nil
		}

		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	cache, err := appConfig.GetAutocertCache()

	if err != nil {
		panic(err)
	}

	certManager := &autocert.Manager{
		Email:      appConfig.Autocert.Email,
		Cache:      cache,
		Client:     &acme.Client{DirectoryURL: appConfig.Autocert.DirectoryURL},
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
	}

	httpsServer := GetWebServer(router, appConfig.Server.HttpsHost)
	httpsServer.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}

	go func() {
		fmt.Printf("HTTPS Server Running: %s\n", httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")

		if err != nil {
			panic(err)
		}
	}()

	httpServer := GetWebServer(GetDefaultWebRouter(), appConfig.Server.HttpHost)
	httpServer.Handler = certManager.HTTPHandler(httpServer.Handler)

	fmt.Printf("HTTP Server Running %s\n", httpServer.Addr)
	err = httpServer.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
