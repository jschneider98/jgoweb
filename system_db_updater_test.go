// +build integration

package jgoweb

import (
	"testing"
)

// ***
func GetTestSystemDbUpdate() *SystemDbUpdate {
	update := CreateSystemDbUpdateNoContext("Test SystemDbUpdate", "Test SystemDbUpdate")

	update.ApplyUpdate = func(ctx ContextInterface) error {
		return nil
	}

	return update
}

//
func TestSystemDbUpdaterGetDbUpdateInfo(t *testing.T) {
	updates := make([]SystemDbUpdateInterface, 0)
	updates = append(updates, GetTestSystemDbUpdate())

	sdu := NewSystemDbUpdater(MockCtx.Db, updates, true)
	sdu.SetDebug(true)

	_, err := sdu.GetDbUpdateInfo()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}
}

//
func TestSystemDbUpdaterRunAll(t *testing.T) {
	updates := make([]SystemDbUpdateInterface, 0)
	updates = append(updates, GetTestSystemDbUpdate())

	sdu := NewSystemDbUpdater(MockCtx.Db, updates, true)
	sdu.SetDebug(true)

	err := sdu.RunAll()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}
}
