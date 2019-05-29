package jgoweb

import (
	"fmt"
	"strings"
	"runtime"
	"io"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"testing"
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
	request, _ := http.NewRequest(method, path, nil)
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
