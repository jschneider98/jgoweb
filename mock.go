package jgoweb

import (
	"github.com/jschneider98/jgoweb/db"
)

var MockDb *db.Collection
var MockUser *User
var MockCtx *WebContext

func InitMockDb() {
	var err error

	if MockDb == nil {
		MockDb, err = db.NewDb()

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockUser() {
	InitMockDb()

	if MockUser == nil {
		var err error

		ctx := NewContext(MockDb)
		MockUser, err = FetchUserByShardEmail(ctx, "jschneider98@gmail.com")

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockCtx() {
	InitMockDb()
	var err error

	if MockCtx == nil {
		MockCtx = &WebContext{}
		MockCtx.Db = MockDb
		MockCtx.DbSess, err = MockDb.GetSessionByName("uxt_0000")

		if err != nil {
			panic(err)
		}
	}
}
