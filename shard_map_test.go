// +build integration

package jgoweb

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"strings"
	"testing"
)

//
func TestFetchShardMapById(t *testing.T) {
	InitMockCtx()
	InitMockShardMap()

	// force not found
	id := "0"
	sm, err := FetchShardMapById(MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch ShardMap by id: %v\n", err)
		return
	}

	if sm != nil {
		t.Errorf("\nERROR: Should have failed to find ShardMap: %v\n", id)
		return
	}

	sm, err = FetchShardMapById(MockCtx, MockShardMap.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if sm == nil {
		t.Errorf("\nERROR: Should have found ShardMap with Id: %v\n", MockShardMap.GetId())
		return
	}

	if sm.GetId() != MockShardMap.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %v Got: %v\n", MockShardMap.GetId(), sm.GetId())
		return
	}
}

//
func TestShardMapId(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetId()
	testVal := "test"

	MockShardMap.SetId("")

	if MockShardMap.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockShardMap.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockShardMap.GetId())
	}

	MockShardMap.SetId(testVal)

	if !MockShardMap.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockShardMap.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockShardMap.GetId())
	}

	MockShardMap.SetId(origVal)
}

//
func TestShardMapShardId(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetShardId()
	testVal := "test"

	MockShardMap.SetShardId("")

	if MockShardMap.ShardId.Valid {
		t.Errorf("ERROR: ShardId should be invalid.\n")
	}

	if MockShardMap.GetShardId() != "" {
		t.Errorf("ERROR: Set ShardId failed. Should have a blank value. Got: %s", MockShardMap.GetShardId())
	}

	MockShardMap.SetShardId(testVal)

	if !MockShardMap.ShardId.Valid {
		t.Errorf("ERROR: ShardId should be valid.\n")
	}

	if MockShardMap.GetShardId() != testVal {
		t.Errorf("ERROR: Set ShardId failed. Expected: %s, Got: %s", testVal, MockShardMap.GetShardId())
	}

	MockShardMap.SetShardId(origVal)
}

//
func TestShardMapDomain(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetDomain()
	testVal := "test"

	MockShardMap.SetDomain("")

	if MockShardMap.Domain.Valid {
		t.Errorf("ERROR: Domain should be invalid.\n")
	}

	if MockShardMap.GetDomain() != "" {
		t.Errorf("ERROR: Set Domain failed. Should have a blank value. Got: %s", MockShardMap.GetDomain())
	}

	MockShardMap.SetDomain(testVal)

	if !MockShardMap.Domain.Valid {
		t.Errorf("ERROR: Domain should be valid.\n")
	}

	if MockShardMap.GetDomain() != testVal {
		t.Errorf("ERROR: Set Domain failed. Expected: %s, Got: %s", testVal, MockShardMap.GetDomain())
	}

	MockShardMap.SetDomain(origVal)
}

//
func TestShardMapAccountId(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetAccountId()
	testVal := "test"

	MockShardMap.SetAccountId("")

	if MockShardMap.AccountId.Valid {
		t.Errorf("ERROR: AccountId should be invalid.\n")
	}

	if MockShardMap.GetAccountId() != "" {
		t.Errorf("ERROR: Set AccountId failed. Should have a blank value. Got: %s", MockShardMap.GetAccountId())
	}

	MockShardMap.SetAccountId(testVal)

	if !MockShardMap.AccountId.Valid {
		t.Errorf("ERROR: AccountId should be valid.\n")
	}

	if MockShardMap.GetAccountId() != testVal {
		t.Errorf("ERROR: Set AccountId failed. Expected: %s, Got: %s", testVal, MockShardMap.GetAccountId())
	}

	MockShardMap.SetAccountId(origVal)
}

//
func TestShardMapCreatedAt(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetCreatedAt()
	testVal := "test"

	MockShardMap.SetCreatedAt("")

	if MockShardMap.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be invalid.\n")
	}

	if MockShardMap.GetCreatedAt() != "" {
		t.Errorf("ERROR: Set CreatedAt failed. Should have a blank value. Got: %s", MockShardMap.GetCreatedAt())
	}

	MockShardMap.SetCreatedAt(testVal)

	if !MockShardMap.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be valid.\n")
	}

	if MockShardMap.GetCreatedAt() != testVal {
		t.Errorf("ERROR: Set CreatedAt failed. Expected: %s, Got: %s", testVal, MockShardMap.GetCreatedAt())
	}

	MockShardMap.SetCreatedAt(origVal)
}

//
func TestShardMapUpdatedAt(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetUpdatedAt()
	testVal := "test"

	MockShardMap.SetUpdatedAt("")

	if MockShardMap.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be invalid.\n")
	}

	if MockShardMap.GetUpdatedAt() != "" {
		t.Errorf("ERROR: Set UpdatedAt failed. Should have a blank value. Got: %s", MockShardMap.GetUpdatedAt())
	}

	MockShardMap.SetUpdatedAt(testVal)

	if !MockShardMap.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be valid.\n")
	}

	if MockShardMap.GetUpdatedAt() != testVal {
		t.Errorf("ERROR: Set UpdatedAt failed. Expected: %s, Got: %s", testVal, MockShardMap.GetUpdatedAt())
	}

	MockShardMap.SetUpdatedAt(origVal)
}

//
func TestShardMapDeletedAt(t *testing.T) {
	InitMockShardMap()
	origVal := MockShardMap.GetDeletedAt()
	testVal := "test"

	MockShardMap.SetDeletedAt("")

	if MockShardMap.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be invalid.\n")
	}

	if MockShardMap.GetDeletedAt() != "" {
		t.Errorf("ERROR: Set DeletedAt failed. Should have a blank value. Got: %s", MockShardMap.GetDeletedAt())
	}

	MockShardMap.SetDeletedAt(testVal)

	if !MockShardMap.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be valid.\n")
	}

	if MockShardMap.GetDeletedAt() != testVal {
		t.Errorf("ERROR: Set DeletedAt failed. Expected: %s, Got: %s", testVal, MockShardMap.GetDeletedAt())
	}

	MockShardMap.SetDeletedAt(origVal)
}

//
func TestShardMapInsert(t *testing.T) {
	InitMockShardMap()
	ShardId := MockShard.GetId()
	Domain := "Domain Insert"
	AccountId := MockUser.GetAccountId()

	sm, err := NewShardMap(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: NewShardMap() failed. %v\n", err)
	}

	sm.SetShardId(ShardId)
	sm.SetDomain(Domain)
	sm.SetAccountId(AccountId)

	err = sm.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !sm.Id.Valid {
		t.Errorf("\nERROR: ShardMap.Id should be set.\n")
	}

	// verify write
	sm, err = FetchShardMapById(MockCtx, sm.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sm == nil || sm.GetShardId() != ShardId || sm.GetDomain() != Domain || sm.GetAccountId() != AccountId {
		t.Errorf("\nERROR: ShardMap does not match save values. Insert failed.\n")
	}
}

//
func TestShardMapUpdate(t *testing.T) {
	InitMockShardMap()
	Domain := "Domain Update"

	MockShardMap.SetDomain(Domain)

	err := MockShardMap.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sm, err := FetchShardMapById(MockCtx, MockShardMap.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sm == nil || sm.GetDomain() != Domain {
		t.Errorf("\nERROR: ShardMap does not match save values. Update failed.\n")
	}
}

//
func TestShardMapDelete(t *testing.T) {
	InitMockShardMap()
	err := MockShardMap.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sm, err := FetchShardMapById(MockCtx, MockShardMap.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !sm.DeletedAt.Valid {
		t.Errorf("\nERROR: ShardMap does not match save values. Delete failed.\n")
	}
}

//
func TestShardMapUndelete(t *testing.T) {
	InitMockShardMap()
	err := MockShardMap.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	err = MockShardMap.Undelete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sm, err := FetchShardMapById(MockCtx, MockShardMap.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sm == nil || sm.DeletedAt.Valid {
		t.Errorf("\nERROR: ShardMap does not match save values. Undelete failed.\n")
	}
}

//
func TestNewShardMapWithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = NewShardMapWithData(MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestShardMapProcessSubmit(t *testing.T) {
	sm, err := NewShardMap(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest(
		"POST",
		"http://example.com",
		strings.NewReader(
			fmt.Sprintf(
				"z=post&ShardId=%s&Domain=Domain&AccountId=%s", MockShard.GetId(), MockUser.GetAccountId(),
			),
		),
	)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := sm.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %v", msg)
	}
}

//
func TestFetchShardMapByAccountId(t *testing.T) {
	_, err := FetchShardMapByAccountId(MockCtx, MockUser.GetAccountId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestGetAllShardMaps(t *testing.T) {
	_, err := GetAllShardMaps(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}
