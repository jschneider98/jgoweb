// +build integration

package jgoweb

import (
	"github.com/gocraft/web"
	"net/http"
	"strings"
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
func TestAccountId(t *testing.T) {
	InitMockAccount()
	origVal := MockAccount.GetId()
	testVal := "test"

	MockAccount.SetId("")

	if MockAccount.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockAccount.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockAccount.GetId())
	}

	MockAccount.SetId(testVal)

	if !MockAccount.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockAccount.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockAccount.GetId())
	}

	MockAccount.SetId(origVal)
}

//
func TestAccountDomain(t *testing.T) {
	InitMockAccount()
	origVal := MockAccount.GetDomain()
	testVal := "test"

	MockAccount.SetDomain("")

	if MockAccount.Domain.Valid {
		t.Errorf("ERROR: Domain should be invalid.\n")
	}

	if MockAccount.GetDomain() != "" {
		t.Errorf("ERROR: Set Domain failed. Should have a blank value. Got: %s", MockAccount.GetDomain())
	}

	MockAccount.SetDomain(testVal)

	if !MockAccount.Domain.Valid {
		t.Errorf("ERROR: Domain should be valid.\n")
	}

	if MockAccount.GetDomain() != testVal {
		t.Errorf("ERROR: Set Domain failed. Expected: %s, Got: %s", testVal, MockAccount.GetDomain())
	}

	MockAccount.SetDomain(origVal)
}

//
func TestAccountCreatedAt(t *testing.T) {
	InitMockAccount()
	origVal := MockAccount.GetCreatedAt()
	testVal := "test"

	MockAccount.SetCreatedAt("")

	if MockAccount.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be invalid.\n")
	}

	if MockAccount.GetCreatedAt() != "" {
		t.Errorf("ERROR: Set CreatedAt failed. Should have a blank value. Got: %s", MockAccount.GetCreatedAt())
	}

	MockAccount.SetCreatedAt(testVal)

	if !MockAccount.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be valid.\n")
	}

	if MockAccount.GetCreatedAt() != testVal {
		t.Errorf("ERROR: Set CreatedAt failed. Expected: %s, Got: %s", testVal, MockAccount.GetCreatedAt())
	}

	MockAccount.SetCreatedAt(origVal)
}

//
func TestAccountUpdatedAt(t *testing.T) {
	InitMockAccount()
	origVal := MockAccount.GetUpdatedAt()
	testVal := "test"

	MockAccount.SetUpdatedAt("")

	if MockAccount.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be invalid.\n")
	}

	if MockAccount.GetUpdatedAt() != "" {
		t.Errorf("ERROR: Set UpdatedAt failed. Should have a blank value. Got: %s", MockAccount.GetUpdatedAt())
	}

	MockAccount.SetUpdatedAt(testVal)

	if !MockAccount.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be valid.\n")
	}

	if MockAccount.GetUpdatedAt() != testVal {
		t.Errorf("ERROR: Set UpdatedAt failed. Expected: %s, Got: %s", testVal, MockAccount.GetUpdatedAt())
	}

	MockAccount.SetUpdatedAt(origVal)
}

//
func TestAccountDeletedAt(t *testing.T) {
	InitMockAccount()
	origVal := MockAccount.GetDeletedAt()
	testVal := "test"

	MockAccount.SetDeletedAt("")

	if MockAccount.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be invalid.\n")
	}

	if MockAccount.GetDeletedAt() != "" {
		t.Errorf("ERROR: Set DeletedAt failed. Should have a blank value. Got: %s", MockAccount.GetDeletedAt())
	}

	MockAccount.SetDeletedAt(testVal)

	if !MockAccount.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be valid.\n")
	}

	if MockAccount.GetDeletedAt() != testVal {
		t.Errorf("ERROR: Set DeletedAt failed. Expected: %s, Got: %s", testVal, MockAccount.GetDeletedAt())
	}

	MockAccount.SetDeletedAt(origVal)
}

//
func TestAccountInsert(t *testing.T) {
	InitMockAccount()
	Domain := "Domain Insert"

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

//
func TestAccountUpdate(t *testing.T) {
	InitMockAccount()
	Domain := "Domain Update"

	MockAccount.SetDomain(Domain)

	err := MockAccount.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	a, err := FetchAccountById(MockCtx, MockAccount.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if a == nil || a.GetDomain() != Domain {
		t.Errorf("\nERROR: Account does not match save values. Update failed.\n")
	}
}

//
func TestAccountDelete(t *testing.T) {
	InitMockAccount()
	err := MockAccount.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	a, err := FetchAccountById(MockCtx, MockAccount.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !a.DeletedAt.Valid {
		t.Errorf("\nERROR: Account does not match save values. Delete failed.\n")
	}
}

//
func TestAccountUndelete(t *testing.T) {
	InitMockAccount()
	err := MockAccount.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	err = MockAccount.Undelete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	a, err := FetchAccountById(MockCtx, MockAccount.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if a == nil || a.DeletedAt.Valid {
		t.Errorf("\nERROR: Account does not match save values. Undelete failed.\n")
	}
}

//
func TestNewAccountWithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = NewAccountWithData(MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestAccountProcessSubmit(t *testing.T) {
	a, err := NewAccount(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", "http://example.com", strings.NewReader("z=post&Domain=Domain"))

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := a.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %v", msg)
	}
}

// ******

//
func TestGetAllAccounts(t *testing.T) {
	InitMockCtx()

	_, err := GetAllAccounts(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestClusterGetAccounts(t *testing.T) {
	InitMockCtx()

	_, err := ClusterGetAccounts(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}
