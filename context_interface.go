package jgoweb

import(
	"github.com/gocraft/dbr"
	jgoWebDb "github.com/jschneider98/jgoweb/db"
	"github.com/gocraft/web"
)

type ContextInterface interface {
	Begin() (*dbr.Tx, error)
	Commit() (error)
	Rollback() (error)
	Select(column ...string) *dbr.SelectBuilder
	SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder
	OptionalBegin() (*dbr.Tx, error)
	OptionalCommit(tx *dbr.Tx) error
	SetUser(user *User)
	SessionGetString(key string) (string, error)
	SessionPutString(rw web.ResponseWriter, key string, value string)
	GetDb() *jgoWebDb.Collection
	GetDbSession() *dbr.Session
	SetDbSession(dbSess *dbr.Session)
}
