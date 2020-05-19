// +build integration

package generator

import (
	"fmt"
	"github.com/jschneider98/jgoweb"
	"os"
	"testing"
)

//
func TestMain(m *testing.M) {
	InitTest()
	code := m.Run()
	TeardownTest()

	os.Exit(code)
}

//
func InitTest() {
	jgoweb.SetConfigEnvVar("JGOWEB_TEST_CONFIG")
	jgoweb.InitMockCtx()

	_, err := jgoweb.MockCtx.Begin()

	if err != nil {
		panic(err)
	}
}

//
func TeardownTest() {
	fmt.Println("Teardown")
	err := jgoweb.MockCtx.Rollback()

	if err != nil {
		panic(err)
	}
}
