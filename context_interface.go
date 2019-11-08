package jgoweb

import(
	"github.com/gocraft/dbr"
	jgoWebDb "github.com/jschneider98/jgoweb/db"
	"gopkg.in/go-playground/validator.v9"
	"github.com/gocraft/web"
)

type ContextInterface interface {
	Begin() (*dbr.Tx, error)
	Commit() (error)
	Rollback() (error)
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
}
