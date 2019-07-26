package jgoweb

import (
	"fmt"
	"time"
	"os"
	// "context"
	"net/http"
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

	sessionManager = scs.NewCookieManager(AppConfig.ServerOptions.SessionKey)
	scs.CookieName = AppConfig.ServerOptions.SessionName
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
	StartHealthSink(AppConfig.ServerOptions.HealthHost)

	if AppConfig.ServerOptions.EnableSsl {
		// @TEMP
		StartHttpServer(router, AppConfig.ServerOptions.HttpHost)
	} else {
		StartHttpServer(router, AppConfig.ServerOptions.HttpHost)
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
