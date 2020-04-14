package jgoweb

import (
	"github.com/gocraft/dbr"
	"github.com/gocraft/web"
	jgoWebDb "github.com/jschneider98/jgoweb/db"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
)

type ContextInterface interface {
	Begin() (*dbr.Tx, error)
	Commit() error
	Rollback() error
	Select(column ...string) *dbr.SelectBuilder
	SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder
	InsertBySql(query string, value ...interface{}) *dbr.InsertStmt
	UpdateBySql(query string, value ...interface{}) *dbr.UpdateStmt
	OptionalBegin() (*dbr.Tx, error)
	OptionalCommit(tx *dbr.Tx) error
	SetUser(user *User)
	SessionGetString(key string) (string, error)
	SessionPutString(rw web.ResponseWriter, key string, value string)
	GetDb() *jgoWebDb.Collection
	GetDbSession() *dbr.Session
	SetDbSession(dbSess *dbr.Session)
	GetValidator() *validator.Validate
	GetTemplate(filename string) (*template.Template, error)
}
