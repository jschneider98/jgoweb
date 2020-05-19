// +build integration

package jgoweb

import (
	"testing"
)

//
func TestTransactions(t *testing.T) {
	// 	InitMockCtx()

	// 	_, err := MockCtx.Begin()

	// 	if err != nil {
	// 		t.Errorf("\nERROR: Failed to start transaction %v\n", err)
	// 	}

	// 	err = MockCtx.Rollback()

	// 	if err != nil {
	// 		t.Errorf("\nERROR: Failed to rollback transaction %v\n", err)
	// 	}

	// 	_, err = MockCtx.Begin()

	// 	if err != nil {
	// 		t.Errorf("\nERROR: Failed to start transaction %v\n", err)
	// 	}

	// 	err = MockCtx.Commit()

	// 	if err != nil {
	// 		t.Errorf("\nERROR: Failed to commit transaction %v\n", err)
	// 	}

	// 	dbSess := MockCtx.GetDbSession()
	// 	dbSess.SelectBySql("SELECT 1")
}
