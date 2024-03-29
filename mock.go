package jgoweb

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

// These aren't true mocks. More like "test" structs
var MockUser *User
var MockCtx *WebContext
var MockAccount *Account
var MockShard *Shard
var MockShardMap *ShardMap
var MockQueueJob *QueueJob
var MockSystemDbUpdate *SystemDbUpdate

//
func InitMockCtx() {
	InitConfig()
	InitDbCollection()
	var err error

	if MockCtx == nil {
		MockCtx = &WebContext{}
		MockCtx.Db = db
		MockCtx.DbSess, err = db.GetSessionByName(appConfig.Integration.ShardName)

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockUser() {
	InitDbCollection()
	InitMockCtx()
	InitConfig()

	if MockUser == nil {
		var err error

		MockUser, err = FetchUserByShardEmail(MockCtx, appConfig.Integration.UserEmail)

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockAccount() {
	InitMockUser()

	if MockAccount == nil {
		var err error

		MockAccount, err = FetchAccountById(MockCtx, MockUser.GetAccountId())

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockQueueJob() {
	InitMockCtx()
	InitMockUser()

	if MockQueueJob == nil {
		var err error

		MockQueueJob, err = NewQueueJob(MockCtx)

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockShard() {
	InitMockUser()

	if MockShard == nil {
		var err error

		MockShard, err = FetchShardByAccountId(MockCtx, MockUser.GetAccountId())

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockShardMap() {
	InitMockUser()

	if MockShardMap == nil {
		var err error

		MockShardMap, err = CreateShardMap(MockCtx,
			MockShard.GetId(),
			"integration_test@test.com",
			MockUser.GetAccountId())

		if err != nil {
			panic(err)
		}

		err = MockShardMap.Save()

		if err != nil {
			panic(err)
		}
	}
}

//
func InitMockSystemDbUpdate() {

	if MockSystemDbUpdate == nil {
		var err error

		MockSystemDbUpdate, err = NewSystemDbUpdate(MockCtx)

		if err != nil {
			panic(err)
		}

		MockSystemDbUpdate.SetUpdateName("Integration Test")
		MockSystemDbUpdate.SetDescription("Integration Test")

		MockSystemDbUpdate.Save()
	}
}

// Return's the caller's caller info.
func CallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

// Make a testing request (lifted/modified from gocraft/web)
func NewTestRequest(method, path string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, body)

	if method == "POST" && body != nil {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	recorder := httptest.NewRecorder()

	return recorder, request
}

//
func AssertResponse(t *testing.T, rr *httptest.ResponseRecorder, code int) {

	if code != rr.Code {
		body, _ := ioutil.ReadAll(rr.Body)

		t.Errorf("assertResponse: expected code to be %d but got %d. (caller: %s) Body: %s", code, rr.Code, CallerInfo(), body)
	}
}
