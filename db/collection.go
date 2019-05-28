package db

import(
	"time"
	"fmt"
	"io/ioutil"
	"math/rand"
	"encoding/json"
	"errors"
	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
)

type CollectionInterface interface {
	GetCons() error
	GetConn(name string) (*dbr.Connection, error)
}

type Collection struct {
	Config []ConnInfo
	Conns map[string]*dbr.Connection
}

// Connection Strings
type ConnInfo struct {
	ShardName      string `json:"shard_name"`
	DbConnString   string `json:"db_conn_string"`
}

// Retrieve db obj
var NewDb = func() (*Collection, error) {
	conns := make(map[string]*dbr.Connection)

	// Load connection string info
	file, err := ioutil.ReadFile("./conns.json")
	
	if err != nil {
		return nil, err
	}

	var connConfig []ConnInfo
	err = json.Unmarshal(file, &connConfig)

	if err != nil {
		return nil, err
	}

	for _, connInfo := range connConfig {
		conn, err := dbr.Open("postgres", connInfo.DbConnString, nil)
		//conn.SetMaxOpenConns(10)

		if err != nil {
			return nil, err
		}

		conns[connInfo.ShardName] = conn;
	}

	db := &Collection{Conns: conns, Config: connConfig}

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

// get Db session by name
func(db *Collection) GetSessionByName(name string) (*dbr.Session, error) {
	
	dbConn, err := db.GetConnByName(name)

	if err != nil {
		return nil, err
	}

	dbSess := dbConn.NewSession(nil)

	return dbSess, nil
}
