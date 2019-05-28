package jgoweb

import(
	"github.com/gocraft/dbr"
	"github.com/jschneider98/jgoweb/db"
)

type ContextInterface interface {
	Begin() (*dbr.Tx, error)
	Commit() (error)
	Rollback() (error)
	Select(column ...string) *dbr.SelectBuilder
	SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder
	OptionalBegin() (*dbr.Tx, error)
	OptionalCommit(tx *dbr.Tx) error
	SetUser(user UserInterface)
	SessionGetString(key string) (string, error)
	GetDb() *db.Collection
	GetDbSession() *dbr.Session
	SetDbSession(dbSess *dbr.Session)
}
