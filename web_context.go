package jgoweb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexedwards/scs"
	"github.com/gocraft/dbr"
	"github.com/gocraft/health"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgovalidator"
	jgoWebDb "github.com/jschneider98/jgoweb/db"
	"github.com/jschneider98/jgoweb/util"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var db *jgoWebDb.Collection

type WebContext struct {
	User                *User
	Session             *scs.Session
	Template            *template.Template
	Job                 *health.Job
	Method              string
	StartTime           time.Time
	EndPoint            string
	Validate            *validator.Validate
	Db                  *jgoWebDb.Collection
	DbSess              *dbr.Session
	Tx                  *dbr.Tx
	WebReqHistogram     *prometheus.HistogramVec
	RollbackTransaction bool
}

var (
	webReqHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "web_request_duration_milliseconds",
		Help: "Histogram of web requests. Labels: method, handler, code.",
	},
		[]string{"method", "handler", "code"},
	)
)

// Init Db
func InitDbCollection() {
	var err error

	InitConfig()

	if db != nil {
		return
	}

	db, err = jgoWebDb.NewDb(appConfig.DbConns)

	if err != nil {
		panic(err)
	}
}

//
func NewContext(db *jgoWebDb.Collection) *WebContext {
	return &WebContext{Db: db, Validate: jgovalidator.GetValidator()}
}

// *** Getters/Setters ***

func (ctx *WebContext) SetUser(user *User) {
	ctx.User = user
}

//
func (ctx *WebContext) SessionGetString(key string) (string, error) {
	return ctx.Session.GetString(key)
}

//
func (ctx *WebContext) SessionPutString(rw web.ResponseWriter, key string, value string) {
	ctx.Session.PutString(rw, key, value)
}

//
func (ctx *WebContext) GetDb() *jgoWebDb.Collection {
	return ctx.Db
}

//
func (ctx *WebContext) GetDbSession() *dbr.Session {
	return ctx.DbSess
}

//
func (ctx *WebContext) GetValidator() *validator.Validate {

	if ctx.Validate == nil {
		ctx.Validate = jgovalidator.GetValidator()
	}

	return ctx.Validate
}

//
func (ctx *WebContext) SetDbSession(dbSess *dbr.Session) {
	ctx.DbSess = dbSess
}

// ******* Db Methods *******

//
func (ctx *WebContext) Begin() (*dbr.Tx, error) {
	var err error

	ctx.Tx, err = ctx.DbSess.Begin()

	return ctx.Tx, err
}

//
func (ctx *WebContext) Commit() error {

	if ctx.Tx == nil {
		return errors.New("Cannot commit. No transaction set in context.")
	}

	err := ctx.Tx.Commit()
	ctx.Tx = nil

	return err
}

//
func (ctx *WebContext) Rollback() error {

	if ctx.Tx == nil {
		return errors.New("Cannot rollback. No transaction set in context.")
	}

	err := ctx.Tx.Rollback()
	ctx.Tx = nil

	return err
}

//
func (ctx *WebContext) Select(column ...string) *dbr.SelectBuilder {
	var stmt *dbr.SelectBuilder

	if ctx.Tx != nil {
		stmt = ctx.Tx.Select(column...)
	} else {
		stmt = ctx.DbSess.Select(column...)
	}

	return stmt
}

//
func (ctx *WebContext) SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder {
	var stmt *dbr.SelectBuilder

	if ctx.Tx != nil {
		stmt = ctx.Tx.SelectBySql(query, value...)
	} else {
		stmt = ctx.DbSess.SelectBySql(query, value...)
	}

	return stmt
}

//
func (ctx *WebContext) Prepare(query string) (*sql.Stmt, error) {

	if ctx.Tx != nil {
		return ctx.Tx.Prepare(query)
	} else {
		return ctx.DbSess.Prepare(query)
	}
}

//
func (ctx *WebContext) InsertBySql(query string, value ...interface{}) *dbr.InsertStmt {
	var stmt *dbr.InsertStmt

	if ctx.Tx != nil {
		stmt = ctx.Tx.InsertBySql(query, value...)
	} else {
		stmt = ctx.DbSess.InsertBySql(query, value...)
	}

	return stmt
}

//
func (ctx *WebContext) UpdateBySql(query string, value ...interface{}) *dbr.UpdateStmt {
	var stmt *dbr.UpdateStmt

	if ctx.Tx != nil {
		stmt = ctx.Tx.UpdateBySql(query, value...)
	} else {
		stmt = ctx.DbSess.UpdateBySql(query, value...)
	}

	return stmt
}

//
func (ctx *WebContext) Update(table string) *dbr.UpdateStmt {
	var stmt *dbr.UpdateStmt

	if ctx.Tx != nil {
		stmt = ctx.Tx.Update(table)
	} else {
		stmt = ctx.DbSess.UpdateBySql(table)
	}

	return stmt
}

// Only start a transaction if one hasn't been started yet
func (ctx *WebContext) OptionalBegin() (*dbr.Tx, error) {

	if ctx.Tx != nil {
		return ctx.Tx, nil
	}

	return ctx.DbSess.Begin()
}

// Commit if there's no tx in the context
func (ctx *WebContext) OptionalCommit(tx *dbr.Tx) error {

	if ctx.Tx != nil {
		return nil
	}

	return tx.Commit()
}

// Complete the DB transaction if the web context is managing it.
func (ctx *WebContext) FinishTransaction() error {

	if ctx.Tx == nil {
		return nil
	}

	if ctx.RollbackTransaction {
		return ctx.Rollback()
	}

	return ctx.Commit()
}

// ******* Job Methods *******

// health stream job error
func (ctx *WebContext) JobError(errorTitle string, err error, codeList ...string) {
	var code string

	ctx.Job.EventErr(errorTitle, err)
	ctx.Job.Complete(health.Error)

	if codeList == nil {
		code = "500"
	} else {
		code = codeList[0]
	}

	ctx.UpdateWebMetrics(code)

	hr := "*******************************\n"
	panic(fmt.Sprintf("%sError: %v\n%s", hr, err, hr))
}

// health stream job success
func (ctx *WebContext) JobSuccess(codeList ...string) {
	var code string

	if codeList == nil {
		code = "200"
	} else {
		code = codeList[0]
	}

	ctx.UpdateWebMetrics(code)
	ctx.Job.Complete(health.Success)

	ctx.FinishTransaction()
}

// health stream job warning
func (ctx *WebContext) JobWarning(title string, err error, codeList ...string) {
	var code string

	ctx.Job.EventErr(title, err)
	ctx.Job.Complete(health.Error)

	if codeList == nil {
		code = "500"
	} else {
		code = codeList[0]
	}

	ctx.UpdateWebMetrics(code)
}

//
func (ctx *WebContext) UpdateWebMetrics(code string) {

	if ctx.WebReqHistogram == nil {
		return
	}

	// convert to milliseconds
	duration := float64(time.Since(ctx.StartTime).Nanoseconds()) / 1000000
	ctx.WebReqHistogram.WithLabelValues(ctx.Method, ctx.EndPoint, code).Observe(duration)
}

// write a JSON response
func (ctx *WebContext) JsonResponse(rw web.ResponseWriter, code int, payload string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	fmt.Fprintf(rw, "%v", payload)
}

// write a JSON error response
func (ctx *WebContext) JsonErrorResponse(rw web.ResponseWriter, code int, err error) {
	payload := fmt.Sprintf(`{"error": %v}`, strconv.Quote(err.Error()))

	if code == 0 {
		code = http.StatusInternalServerError
	}

	ctx.JsonResponse(rw, code, payload)
}

// write a JSON error response
func (ctx *WebContext) JsonOkResponse(rw web.ResponseWriter, code int, message string) {
	payload := fmt.Sprintf(`{"message": %v}`, strconv.Quote(message))

	if code == 0 {
		code = http.StatusOK
	}

	ctx.JsonResponse(rw, code, payload)
}

// Init Db
func (ctx *WebContext) InitDbSession() {
	var err error

	if ctx.Db == nil {
		ctx.Db = db
	}

	dbConn, err := ctx.Db.GetRandomConn()

	if err != nil {
		panic(err)
	}

	ctx.DbSess = dbConn.NewSession(nil)
}

// Init Metrics
func (ctx *WebContext) InitMetrics() {
	ctx.WebReqHistogram = webReqHistogram
	prometheus.Register(ctx.WebReqHistogram)
}

// **** Middleware ****

// Auth middleware
func (ctx *WebContext) AjaxRequireUser(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	var err error

	if ctx.User == nil {
		ctx.User, err = NewUser(ctx)

		if err != nil {
			ctx.JobError(util.WhereAmI(), err)
		}
	}

	err = ctx.User.SetFromSession()

	if err != nil {
		fmt.Fprint(rw, "{error: 'User authendication required.'}")
		ctx.JobSuccess()
	} else {
		next(rw, req)
	}
}

// Various context dependancy injection (DbCollection, metrics, etc)
func (ctx *WebContext) LoadDi(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	ctx.InitDbSession()
	ctx.InitMetrics()

	ctx.Method = strings.ToLower(req.Method)
	ctx.StartTime = time.Now()

	next(rw, req)
}

// Have the context manage the DB transaction. Typically don't do this. Each model handles it's own transaction,
// which is much faster/safer than having a single transaction for an entire web request.
// This middleware is useful for testing routes though.
func (ctx *WebContext) BeginTransaction(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	_, err := ctx.Begin()

	if err != nil {
		panic(err)
	}

	next(rw, req)
}

// Tell the web context to rollback it's transaction upon Job success
// This middleware is useful for testing routes.
func (ctx *WebContext) SetJobSuccessRollback(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	ctx.RollbackTransaction = true

	next(rw, req)
}

// EndPoint middleware
// Get endpoint for route (i.e., path = "/search/other/etc", endpoint = "search")
func (ctx *WebContext) LoadEndPoint(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	matches := strings.Split(req.URL.Path, "/")

	if len(matches) >= 2 {
		ctx.EndPoint = matches[1]
	}

	// Special ajax case
	if ctx.EndPoint == "ajax" && len(matches) >= 3 {
		ctx.EndPoint = matches[2]
	}

	next(rw, req)
}

// HealthStream job middle ware for routes
func (ctx *WebContext) LoadJob(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	method := strings.ToLower(req.Method)
	endPoint := ctx.EndPoint

	if method == "" {
		method = "get"
	}

	if endPoint == "" {
		endPoint = "root"
	}

	ctx.Job = healthStream.NewJob("route_" + method + "_" + endPoint)

	next(rw, req)
}

// Session middleware
func (ctx *WebContext) LoadSession(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	ctx.Session = sessionManager.Load(req.Request)
	ctx.Session.RenewToken(rw)
	ctx.Session.PutString(rw, "init", "1")

	next(rw, req)
}

// Template middleware
func (ctx *WebContext) LoadTemplate(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	layout := filepath.Join("static", "templates", "layout.html")
	routeTemplate := ""

	// Conditionally load route template file (if it exists)
	if ctx.EndPoint != "" {
		routeTemplate = filepath.Join("static", "templates", ctx.EndPoint+".html")

		if _, err := os.Stat(routeTemplate); os.IsNotExist(err) {
			routeTemplate = ""
		}
	}

	name := path.Base(layout)
	tmpl, err := template.New(name).Delims("[[", "]]").ParseFiles(layout)

	if err != nil {
		ctx.JobError(util.WhereAmI(), err)
	}

	if routeTemplate != "" {
		tmpl, err = tmpl.ParseFiles(routeTemplate)

		if err != nil {
			ctx.JobError(util.WhereAmI(), err)
		}
	}

	ctx.Template = tmpl

	next(rw, req)
}

// Auth middleware
func (ctx *WebContext) RequireUser(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	var err error

	if ctx.User == nil {
		ctx.User, err = NewUser(ctx)

		if err != nil {
			ctx.JobError(util.WhereAmI(), err)
		}
	}

	err = ctx.User.SetFromSession()

	if err != nil {
		// Force sign in
		http.Redirect(rw, req.Request, "/login", 302)
		return
	}

	next(rw, req)
}

// ******

// Get Template
func (ctx *WebContext) GetTemplate(filename string) (*template.Template, error) {
	file := filepath.Join("static", "templates", filename)

	name := path.Base(file)
	return template.New(name).Delims("[[", "]]").ParseFiles(file)
}
