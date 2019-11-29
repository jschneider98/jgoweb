package db

import(
	"time"
	"fmt"
	"math/rand"
	"errors"
	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
	"github.com/jschneider98/jgoweb/config"
)

type CollectionInterface interface {
	GetConns() map[string]*dbr.Connection
	GetConnByName(name string) (*dbr.Connection, error)
}

type Collection struct {
	Config []config.DbConnOptions
	Conns map[string]*dbr.Connection
}

// Retrieve db obj
var NewDb = func(dbConns []config.DbConnOptions) (*Collection, error) {
	conns := make(map[string]*dbr.Connection)

	for _, connInfo := range dbConns {
		conn, err := dbr.Open("postgres", connInfo.Dsn, nil)
		//conn.SetMaxOpenConns(10)

		if err != nil {
			return nil, err
		}

		conns[connInfo.ShardName] = conn;
	}

	db := &Collection{Conns: conns, Config: dbConns}

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
func(db *Collection) GetSessionByName(name string) (*dbr.Session, error) {
	
	dbConn, err := db.GetConnByName(name)

	if err != nil {
		return nil, err
	}

	dbSess := dbConn.NewSession(nil)

	return dbSess, nil
}
