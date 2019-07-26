package jgoweb

import (
	"fmt"
	"github.com/gocraft/web"
)

func (c *WebContext) Welcome(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "<html><body>Welcome.</body></html>")
}

//
func GetDefaultWebRouter() *web.Router {
	webContext := WebContext{}

	router := web.New(webContext).
		Middleware(web.LoggerMiddleware).
		Get("/", (*WebContext).Welcome)

	return router
}