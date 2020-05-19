// +build integration

package jgoweb

import (
	"github.com/gocraft/web"
	"net/http"
	"strings"
	"testing"
)

//
func TestFetchSystemDbUpdateById(t *testing.T) {
	InitMockCtx()
	InitMockSystemDbUpdate()

	// force not found
	id := "0"
	sdu, err := FetchSystemDbUpdateById(MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch SystemDbUpdate by id: %v\n", err)
		return
	}

	if sdu != nil {
		t.Errorf("\nERROR: Should have failed to find SystemDbUpdate: %v\n", id)
		return
	}

	sdu, err = FetchSystemDbUpdateById(MockCtx, MockSystemDbUpdate.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if sdu == nil {
		t.Errorf("\nERROR: Should have found SystemDbUpdate with Id: %v\n", MockSystemDbUpdate.GetId())
		return
	}

	if sdu.GetId() != MockSystemDbUpdate.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %v Got: %v\n", MockSystemDbUpdate.GetId(), sdu.GetId())
		return
	}
}

//
func TestSystemDbUpdateId(t *testing.T) {
	InitMockSystemDbUpdate()
	origVal := MockSystemDbUpdate.GetId()
	testVal := "test"

	MockSystemDbUpdate.SetId("")

	if MockSystemDbUpdate.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockSystemDbUpdate.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockSystemDbUpdate.GetId())
	}

	MockSystemDbUpdate.SetId(testVal)

	if !MockSystemDbUpdate.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockSystemDbUpdate.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockSystemDbUpdate.GetId())
	}

	MockSystemDbUpdate.SetId(origVal)
}

//
func TestSystemDbUpdateUpdateName(t *testing.T) {
	InitMockSystemDbUpdate()
	origVal := MockSystemDbUpdate.GetUpdateName()
	testVal := "test"

	MockSystemDbUpdate.SetUpdateName("")

	if MockSystemDbUpdate.UpdateName.Valid {
		t.Errorf("ERROR: UpdateName should be invalid.\n")
	}

	if MockSystemDbUpdate.GetUpdateName() != "" {
		t.Errorf("ERROR: Set UpdateName failed. Should have a blank value. Got: %s", MockSystemDbUpdate.GetUpdateName())
	}

	MockSystemDbUpdate.SetUpdateName(testVal)

	if !MockSystemDbUpdate.UpdateName.Valid {
		t.Errorf("ERROR: UpdateName should be valid.\n")
	}

	if MockSystemDbUpdate.GetUpdateName() != testVal {
		t.Errorf("ERROR: Set UpdateName failed. Expected: %s, Got: %s", testVal, MockSystemDbUpdate.GetUpdateName())
	}

	MockSystemDbUpdate.SetUpdateName(origVal)
}

//
func TestSystemDbUpdateDescription(t *testing.T) {
	InitMockSystemDbUpdate()
	origVal := MockSystemDbUpdate.GetDescription()
	testVal := "test"

	MockSystemDbUpdate.SetDescription("")

	if MockSystemDbUpdate.Description.Valid {
		t.Errorf("ERROR: Description should be invalid.\n")
	}

	if MockSystemDbUpdate.GetDescription() != "" {
		t.Errorf("ERROR: Set Description failed. Should have a blank value. Got: %s", MockSystemDbUpdate.GetDescription())
	}

	MockSystemDbUpdate.SetDescription(testVal)

	if !MockSystemDbUpdate.Description.Valid {
		t.Errorf("ERROR: Description should be valid.\n")
	}

	if MockSystemDbUpdate.GetDescription() != testVal {
		t.Errorf("ERROR: Set Description failed. Expected: %s, Got: %s", testVal, MockSystemDbUpdate.GetDescription())
	}

	MockSystemDbUpdate.SetDescription(origVal)
}

//
func TestSystemDbUpdateCreatedAt(t *testing.T) {
	InitMockSystemDbUpdate()
	origVal := MockSystemDbUpdate.GetCreatedAt()
	testVal := "test"

	MockSystemDbUpdate.SetCreatedAt("")

	if MockSystemDbUpdate.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be invalid.\n")
	}

	if MockSystemDbUpdate.GetCreatedAt() != "" {
		t.Errorf("ERROR: Set CreatedAt failed. Should have a blank value. Got: %s", MockSystemDbUpdate.GetCreatedAt())
	}

	MockSystemDbUpdate.SetCreatedAt(testVal)

	if !MockSystemDbUpdate.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be valid.\n")
	}

	if MockSystemDbUpdate.GetCreatedAt() != testVal {
		t.Errorf("ERROR: Set CreatedAt failed. Expected: %s, Got: %s", testVal, MockSystemDbUpdate.GetCreatedAt())
	}

	MockSystemDbUpdate.SetCreatedAt(origVal)
}

//
func TestSystemDbUpdateInsert(t *testing.T) {
	InitMockSystemDbUpdate()
	UpdateName := "UpdateName Insert"
	Description := "Description Insert"

	sdu, err := NewSystemDbUpdate(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: NewSystemDbUpdate() failed. %v\n", err)
	}

	sdu.SetUpdateName(UpdateName)
	sdu.SetDescription(Description)

	err = sdu.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !sdu.Id.Valid {
		t.Errorf("\nERROR: SystemDbUpdate.Id should be set.\n")
	}

	// verify write
	sdu, err = FetchSystemDbUpdateById(MockCtx, sdu.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sdu == nil || sdu.GetUpdateName() != UpdateName || sdu.GetDescription() != Description {
		t.Errorf("\nERROR: SystemDbUpdate does not match save values. Insert failed.\n")
	}
}

//
func TestSystemDbUpdateUpdate(t *testing.T) {
	InitMockSystemDbUpdate()
	UpdateName := "UpdateName Update"
	Description := "Description Update"

	MockSystemDbUpdate.SetUpdateName(UpdateName)
	MockSystemDbUpdate.SetDescription(Description)

	err := MockSystemDbUpdate.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sdu, err := FetchSystemDbUpdateById(MockCtx, MockSystemDbUpdate.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sdu == nil || sdu.GetUpdateName() != UpdateName || sdu.GetDescription() != Description {
		t.Errorf("\nERROR: SystemDbUpdate does not match save values. Update failed.\n")
	}
}

//
func TestSystemDbUpdateDelete(t *testing.T) {
	InitMockSystemDbUpdate()
	err := MockSystemDbUpdate.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sdu, err := FetchSystemDbUpdateById(MockCtx, MockSystemDbUpdate.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sdu != nil {
		t.Errorf("\nERROR: Delete failed. Fetch should return nil.\n")
		return
	}

	MockSystemDbUpdate = nil
}

//
func TestNewSystemDbUpdateWithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = NewSystemDbUpdateWithData(MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestSystemDbUpdateProcessSubmit(t *testing.T) {
	sdu, err := NewSystemDbUpdate(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", "http://example.com", strings.NewReader("z=post&UpdateName=UpdateName&Description=Description"))

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := sdu.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %v", msg)
	}
}

//
func TestCreateSystemDbUpdateByUpdateName(t *testing.T) {
	_, err := CreateSystemDbUpdateByUpdateName(MockCtx, "New test update")

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestSystemDbUpdateRun(t *testing.T) {
	sdu := CreateSystemDbUpdateNoContext("New test update", "New test update")
	sdu.SetContext(MockCtx)
	sdu.Clone()

	sdu.ApplyUpdate = func(ctx ContextInterface) error {
		return nil
	}

	needsToRun, err := sdu.NeedsToRun()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !needsToRun {
		t.Errorf("\nERROR: DB update should need to run.\n")
	}

	err = sdu.Run()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	err = sdu.SetComplete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}
}
