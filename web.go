package jgoweb

import (
	"fmt"
	"time"
	"os"
	"github.com/gocraft/health"
	"github.com/alexedwards/scs"
)

var healthStream = health.NewStream()

// @TODO: pull in secret via file/environment
var sessionManager = scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE3")

// Init session
func InitSession(sessionName string) {
	scs.CookieName = sessionName
}

// Start Health Sink
func StartHealthSink(hostname string) {
	healthStream.AddSink(&health.WriterSink{os.Stdout})
	sink := health.NewJsonPollingSink(time.Minute, time.Minute*5)
	healthStream.AddSink(sink)

	fmt.Println("Health Sink Running: ", hostname)
	sink.StartServer(hostname)
}
