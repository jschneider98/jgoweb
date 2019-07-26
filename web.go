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
var AppConfig *config.Config

// Init config
func InitConfig(path string) {
	var err error

	AppConfig, err = config.New(path)

	if err != nil {
		panic(err)
	}
}

// Init session
func InitSession() {

	if AppConfig == nil {
		InitConfig("./config/config.json")
	}

	sessionManager = scs.NewCookieManager(AppConfig.Server.SessionKey)
	scs.CookieName = AppConfig.Server.SessionName
}

//
func Start(router *web.Router) {
	InitConfig("./config/config.json")
	StartAll(router)
}

//
func StartAll(router *web.Router) {

	if AppConfig == nil {
		InitConfig("./config/config.json")
	}

	InitDbCollection()
	InitSession()
	StartHealthSink(AppConfig.Server.HealthHost)

	if AppConfig.Server.EnableSsl {
		StartHttpsServer(router)
	} else {
		StartHttpServer(router, AppConfig.Server.HttpHost)
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

	if AppConfig == nil {
		InitConfig("./config/config.json")
	}

	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := AppConfig.Autocert.AllowedHost

		if host == allowedHost {
			return nil
		}

		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	cache, err := AppConfig.GetAutocertCache()

	if err != nil {
		panic(err)
	}

	acm := &autocert.Manager{
		Email:      AppConfig.Autocert.Email,
		Cache:      cache,
		Client:     &acme.Client{DirectoryURL: AppConfig.Autocert.DirectoryURL},
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
	}

	httpsServer := GetWebServer(router, AppConfig.Server.HttpsHost)
	httpsServer.TLSConfig = &tls.Config{GetCertificate: acm.GetCertificate}

	go func() {
		fmt.Printf("HTTPS Server Running: %s\n", httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")

		if err != nil {
			panic(err)
		}
	}()
}
