package db

import (
	"errors"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/jschneider98/jgoweb/config"
	_ "github.com/lib/pq"
	"math/rand"
	"time"
)

type CollectionInterface interface {
	GetConns() map[string]*dbr.Connection
	GetConnByName(name string) (*dbr.Connection, error)
}

type Collection struct {
	Config    []config.DbConnOptions
	ConfigMap map[string]int
	Conns     map[string]*dbr.Connection
}

// Retrieve db obj
var NewDb = func(dbConns []config.DbConnOptions) (*Collection, error) {
	conns := make(map[string]*dbr.Connection)
	configMap := make(map[string]int)
	var maxOpenConns int
	var maxIdleConns int
	var connMaxLifetime int

	for index, connInfo := range dbConns {
		conn, err := dbr.Open("postgres", connInfo.Dsn, nil)

		// defaults
		maxOpenConns = 100
		maxIdleConns = 25
		connMaxLifetime = 30

		if connInfo.MaxOpenConns > 0 {
			maxOpenConns = connInfo.MaxOpenConns
		}

		if connInfo.MaxIdleConns > 0 {
			maxIdleConns = connInfo.MaxIdleConns
		}

		if connInfo.ConnMaxLifetime > 0 {
			connMaxLifetime = connInfo.ConnMaxLifetime
		}

		conn.SetMaxOpenConns(maxOpenConns)
		conn.SetMaxIdleConns(maxIdleConns)
		conn.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)

		if err != nil {
			return nil, err
		}

		conns[connInfo.ShardName] = conn
		configMap[connInfo.ShardName] = index
	}

	db := &Collection{Conns: conns, Config: dbConns, ConfigMap: configMap}

	return db, nil
}

// get Db connection by name
func (db *Collection) GetConnByName(name string) (*dbr.Connection, error) {

	if conn, ok := db.Conns[name]; ok {
		return conn, nil
	}

	err := errors.New(fmt.Sprintf("Connection %s does not exist", name))
	return nil, err
}

// get DB Config by name
func (db *Collection) GetConfigByName(name string) (config.DbConnOptions, error) {
	var empty config.DbConnOptions

	if index, ok := db.ConfigMap[name]; ok {
		return db.Config[index], nil
	}

	err := errors.New(fmt.Sprintf("Config %s does not exist", name))
	return empty, err
}

// get random DB conn
func (db *Collection) GetRandomConn() (*dbr.Connection, error) {

	if len(db.Config) == 0 {
		err := errors.New("Empty DB Config")
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	connInfo := db.Config[rand.Intn(len(db.Config))]

	return db.GetConnByName(connInfo.ShardName)
}

//
func (db *Collection) GetConns() map[string]*dbr.Connection {
	return db.Conns
}

// get Db session by name
func (db *Collection) GetSessionByName(name string) (*dbr.Session, error) {

	dbConn, err := db.GetConnByName(name)

	if err != nil {
		return nil, err
	}

	config, err := db.GetConfigByName(name)

	if err != nil {
		return nil, err
	}

	dbSess := dbConn.NewSession(nil)
	dbSess.Timeout = time.Duration(config.StatementTimeout) * time.Millisecond

	return dbSess, nil
}
