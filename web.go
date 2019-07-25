package jgoweb

import (
	"fmt"
	"flag"
	"time"
	"os"
	"net/http"
	"github.com/gocraft/web"
	"github.com/gocraft/health"
	"github.com/alexedwards/scs"
)

type WebParams struct {
	SessionName string `json:sessionName`
	EnableSsl bool `json:enableSsl`
	HttpsHost string `json:httpsHost`
	HttpHost string `json:httpHost`
	HealthHost string `json:healthHost`
}

var webParams WebParams
var healthStream = health.NewStream()

// @TODO: pull in secret via file/environment
var sessionManager = scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE3")

// Init session
func InitSession(sessionName string) {
	scs.CookieName = sessionName
}

//
func ParseFlags() *WebParams {
	flag.StringVar(&webParams.SessionName, "session", "web-session", "Session name")
	flag.StringVar(&webParams.SessionName, "s", "web-session", "Session name")

	flag.BoolVar(&webParams.EnableSsl, "ssl", true, "Enable SSL (true/false)")
	flag.StringVar(&webParams.HttpHost, "http", ":80", "HTTP host with optional port (i.e., localhost:80)")
	flag.StringVar(&webParams.HttpsHost, "https", ":443", "HTTPS host with optional port (i.e., localhost:443)")

	flag.StringVar(&webParams.HealthHost, "health-host", ":5020", "HTTP health sink host with optional port (i.e., localhost:5020)")
	flag.StringVar(&webParams.HealthHost, "hh", ":5020", "HTTP health sink host with optional port (i.e., localhost:5020)")

	flag.Parse()

	return &webParams
}

//
func Start(router *web.Router) {
	ParseFlags()
	StartAll(router, &webParams)
}

//
func StartAll(router *web.Router, params *WebParams) {
	InitDbCollection()
	InitSession(params.SessionName)
	StartHealthSink(params.HealthHost)

	server := GetWebServer(router, params.HttpHost)
	StartHttpServer(server)
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
func StartHttpServer(server *http.Server) {
	fmt.Println("HTTP Server Running: ", server.Addr)
	server.ListenAndServe()
}
