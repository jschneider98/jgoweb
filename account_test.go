// +build integration

package jgoweb

import (
	"testing"
)

//
func TestFetchAccountById(t *testing.T) {
	InitMockCtx()
	InitMockAccount()

	// force not found
	id := "00000000-0000-0000-0000-000000000000"
	a, err := FetchAccountById(MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch Account by id: %v\n", err)
		return
	}

	if a != nil {
		t.Errorf("\nERROR: Should have failed to find Account: %v\n", id)
		return
	}

	a, err = FetchAccountById(MockCtx, MockAccount.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if a == nil {
		t.Errorf("\nERROR: Should have found Account with Id: %v\n", MockAccount.GetId())
		return
	}

	if a.GetId() != MockAccount.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %v Got: %v\n", MockAccount.GetId(), a.GetId())
		return
	}
}

//
func TestAccountInsert(t *testing.T) {
	InitMockAccount()
	Domain := "Domain"

	a, err := NewAccount(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: NewAccount() failed. %v\n", err)
	}

	a.SetDomain(Domain)

	err = a.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !a.Id.Valid {
		t.Errorf("\nERROR: Account.Id should be set.\n")
	}

	// verify write
	a, err = FetchAccountById(MockCtx, a.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if a == nil || a.GetDomain() != Domain {
		t.Errorf("\nERROR: Account does not match save values. Insert failed.\n")
	}
}
