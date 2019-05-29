// +build integration

package jgoweb

import (
	"errors"
	"testing"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
)

func (ctx *WebContext) SetTestUser(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// Test loading user from session
	ctx.SessionPutString(rw, "user_email", "jschneider98@gmail.com")

	next(rw, req)
}

// GET: Index route
func (ctx *WebContext) index(rw web.ResponseWriter, req *web.Request) {
	ctx.JobSuccess()
}

// GET: Index route
func (ctx *WebContext) routeError(rw web.ResponseWriter, req *web.Request) {
	err := errors.New("Forced error")
	ctx.JobError(util.WhereAmI(), err)
}

//
func TestRoutes(t *testing.T) {

	router := web.New(WebContext{}).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*WebContext).LoadDi).
		Middleware((*WebContext).LoadEndPoint).
		Middleware((*WebContext).LoadTemplate).
		Middleware((*WebContext).LoadJob).
		Middleware((*WebContext).LoadSession).
		Get("/index", (*WebContext).index)

	router.Subrouter(WebContext{}, "/").
		Middleware((*WebContext).RequireUser).
		Get("/require_user", (*WebContext).index)

	router.Subrouter(WebContext{}, "/").
		Middleware((*WebContext).AjaxRequireUser).
		Get("/ajax_require_user", (*WebContext).index)

	router.Subrouter(WebContext{}, "/").
		Middleware((*WebContext).SetTestUser).
		Middleware((*WebContext).RequireUser).
		Get("/logged_in_index", (*WebContext).index)

	router.Subrouter(WebContext{}, "/").
		Get("/forced_error", (*WebContext).routeError)

	rw, req := NewTestRequest("GET", "/index", nil)
	router.ServeHTTP(rw, req)
	AssertResponse(t, rw, 200)

	//
	rw, req = NewTestRequest("GET", "/require_user", nil)
	router.ServeHTTP(rw, req)
	AssertResponse(t, rw, 200)

	//
	rw, req = NewTestRequest("GET", "/ajax_require_user", nil)
	router.ServeHTTP(rw, req)
	AssertResponse(t, rw, 200)

	//
	rw, req = NewTestRequest("GET", "/logged_in_index", nil)
	router.ServeHTTP(rw, req)
	AssertResponse(t, rw, 200)

	//
	rw, req = NewTestRequest("GET", "/forced_error", nil)
	router.ServeHTTP(rw, req)
	AssertResponse(t, rw, 500)
}
