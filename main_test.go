// +build integration

package jgoweb

import (
	"fmt"
)

import (
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
	fmt.Println("Init Tests")
	SetConfigEnvVar("JGOWEB_TEST_CONFIG")
	InitMockCtx()
	_, err := MockCtx.Begin()

	if err != nil {
		panic(err)
	}

	InitMockUser()
	InitMockAccount()
	InitMockShard()
	InitMockShardMap()
	InitMockSystemDbUpdate()
}

//
func TeardownTest() {
	fmt.Println("Teardown")
	err := MockCtx.Rollback()

	if err != nil {
		panic(err)
	}
}
