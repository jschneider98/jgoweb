// +build integration

package jgoweb

import (
	"github.com/gocraft/web"
	"net/http"
	"strings"
	"testing"
)

//
func TestFetchShardById(t *testing.T) {
	InitMockCtx()
	InitMockShard()

	// force not found
	id := "0"
	s, err := FetchShardById(MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch Shard by id: %v\n", err)
		return
	}

	if s != nil {
		t.Errorf("\nERROR: Should have failed to find Shard: %v\n", id)
		return
	}

	s, err = FetchShardById(MockCtx, MockShard.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if s == nil {
		t.Errorf("\nERROR: Should have found Shard with Id: %v\n", MockShard.GetId())
		return
	}

	if s.GetId() != MockShard.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %v Got: %v\n", MockShard.GetId(), s.GetId())
		return
	}
}

//
func TestShardId(t *testing.T) {
	InitMockShard()
	origVal := MockShard.GetId()
	testVal := "test"

	MockShard.SetId("")

	if MockShard.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockShard.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockShard.GetId())
	}

	MockShard.SetId(testVal)

	if !MockShard.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockShard.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockShard.GetId())
	}

	MockShard.SetId(origVal)
}

//
func TestShardName(t *testing.T) {
	InitMockShard()
	origVal := MockShard.GetName()
	testVal := "test"

	MockShard.SetName("")

	if MockShard.Name.Valid {
		t.Errorf("ERROR: Name should be invalid.\n")
	}

	if MockShard.GetName() != "" {
		t.Errorf("ERROR: Set Name failed. Should have a blank value. Got: %s", MockShard.GetName())
	}

	MockShard.SetName(testVal)

	if !MockShard.Name.Valid {
		t.Errorf("ERROR: Name should be valid.\n")
	}

	if MockShard.GetName() != testVal {
		t.Errorf("ERROR: Set Name failed. Expected: %s, Got: %s", testVal, MockShard.GetName())
	}

	MockShard.SetName(origVal)
}

//
func TestShardAccountCount(t *testing.T) {
	InitMockShard()
	origVal := MockShard.GetAccountCount()
	testVal := "1"

	MockShard.SetAccountCount("")

	if MockShard.AccountCount.Valid {
		t.Errorf("ERROR: AccountCount should be invalid.\n")
	}

	if MockShard.GetAccountCount() != "" {
		t.Errorf("ERROR: Set AccountCount failed. Should have a blank value. Got: %s", MockShard.GetAccountCount())
	}

	MockShard.SetAccountCount(testVal)

	if !MockShard.AccountCount.Valid {
		t.Errorf("ERROR: AccountCount should be valid.\n")
	}

	if MockShard.GetAccountCount() != testVal {
		t.Errorf("ERROR: Set AccountCount failed. Expected: %s, Got: %s", testVal, MockShard.GetAccountCount())
	}

	MockShard.SetAccountCount(origVal)
}

//
func TestShardCreatedAt(t *testing.T) {
	InitMockShard()
	origVal := MockShard.GetCreatedAt()
	testVal := "test"

	MockShard.SetCreatedAt("")

	if MockShard.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be invalid.\n")
	}

	if MockShard.GetCreatedAt() != "" {
		t.Errorf("ERROR: Set CreatedAt failed. Should have a blank value. Got: %s", MockShard.GetCreatedAt())
	}

	MockShard.SetCreatedAt(testVal)

	if !MockShard.CreatedAt.Valid {
		t.Errorf("ERROR: CreatedAt should be valid.\n")
	}

	if MockShard.GetCreatedAt() != testVal {
		t.Errorf("ERROR: Set CreatedAt failed. Expected: %s, Got: %s", testVal, MockShard.GetCreatedAt())
	}

	MockShard.SetCreatedAt(origVal)
}

//
func TestShardUpdatedAt(t *testing.T) {
	InitMockShard()
	origVal := MockShard.GetUpdatedAt()
	testVal := "test"

	MockShard.SetUpdatedAt("")

	if MockShard.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be invalid.\n")
	}

	if MockShard.GetUpdatedAt() != "" {
		t.Errorf("ERROR: Set UpdatedAt failed. Should have a blank value. Got: %s", MockShard.GetUpdatedAt())
	}

	MockShard.SetUpdatedAt(testVal)

	if !MockShard.UpdatedAt.Valid {
		t.Errorf("ERROR: UpdatedAt should be valid.\n")
	}

	if MockShard.GetUpdatedAt() != testVal {
		t.Errorf("ERROR: Set UpdatedAt failed. Expected: %s, Got: %s", testVal, MockShard.GetUpdatedAt())
	}

	MockShard.SetUpdatedAt(origVal)
}

//
func TestShardDeletedAt(t *testing.T) {
	InitMockShard()
	origVal := MockShard.GetDeletedAt()
	testVal := "test"

	MockShard.SetDeletedAt("")

	if MockShard.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be invalid.\n")
	}

	if MockShard.GetDeletedAt() != "" {
		t.Errorf("ERROR: Set DeletedAt failed. Should have a blank value. Got: %s", MockShard.GetDeletedAt())
	}

	MockShard.SetDeletedAt(testVal)

	if !MockShard.DeletedAt.Valid {
		t.Errorf("ERROR: DeletedAt should be valid.\n")
	}

	if MockShard.GetDeletedAt() != testVal {
		t.Errorf("ERROR: Set DeletedAt failed. Expected: %s, Got: %s", testVal, MockShard.GetDeletedAt())
	}

	MockShard.SetDeletedAt(origVal)
}

//
func TestShardInsert(t *testing.T) {
	InitMockShard()
	Name := "Name Insert"
	AccountCount := "1"

	s, err := NewShard(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: NewShard() failed. %v\n", err)
	}

	s.SetName(Name)
	s.SetAccountCount(AccountCount)

	err = s.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !s.Id.Valid {
		t.Errorf("\nERROR: Shard.Id should be set.\n")
	}

	// verify write
	s, err = FetchShardById(MockCtx, s.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if s == nil || s.GetName() != Name || s.GetAccountCount() != AccountCount {
		t.Errorf("\nERROR: Shard does not match save values. Insert failed.\n")
	}
}

//
func TestShardUpdate(t *testing.T) {
	InitMockShard()
	origName := MockShard.GetName()
	Name := "Name Update"
	AccountCount := "2"

	MockShard.SetName(Name)
	MockShard.SetAccountCount(AccountCount)

	err := MockShard.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	s, err := FetchShardById(MockCtx, MockShard.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if s == nil || s.GetName() != Name || s.GetAccountCount() != AccountCount {
		t.Errorf("\nERROR: Shard does not match save values. Update failed.\n")
	}

	MockShard.SetName(origName)
}

//
func TestShardDelete(t *testing.T) {
	InitMockShard()
	err := MockShard.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	s, err := FetchShardById(MockCtx, MockShard.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !s.DeletedAt.Valid {
		t.Errorf("\nERROR: Shard does not match save values. Delete failed.\n")
	}
}

//
func TestShardUndelete(t *testing.T) {
	InitMockShard()
	err := MockShard.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	err = MockShard.Undelete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	s, err := FetchShardById(MockCtx, MockShard.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if s == nil || s.DeletedAt.Valid {
		t.Errorf("\nERROR: Shard does not match save values. Undelete failed.\n")
	}
}

//
func TestNewShardWithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = NewShardWithData(MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestShardProcessSubmit(t *testing.T) {
	s, err := NewShard(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", "http://example.com", strings.NewReader("z=post&Name=Name&AccountCount=1"))

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := s.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %v", msg)
	}
}

// *****

//
func TestFetchShardByAccountId(t *testing.T) {
	InitMockCtx()
	InitMockUser()

	accountId := "00000000-0000-0000-0000-000000000000"

	shard, err := FetchShardByAccountId(MockCtx, accountId)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if shard != nil {
		t.Errorf("\nInvalid account ID should return nil.\n")
	}

	shard, err = FetchShardByAccountId(MockCtx, MockUser.GetAccountId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if shard == nil {
		t.Errorf("\nValid account ID should return a shard.\n")
	}
}

//
func TestShardNewWebContext(t *testing.T) {
	InitMockShard()

	_, err := MockShard.NewWebContext()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestGetShardByAccountId(t *testing.T) {
	InitMockCtx()
	InitMockUser()

	_, err := GetShardByAccountId(MockCtx, MockUser.GetAccountId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestFetchShardByName(t *testing.T) {
	InitMockCtx()
	InitMockShard()

	_, err := FetchShardByName(MockCtx, MockShard.GetName())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestCreateShardByName(t *testing.T) {
	InitMockCtx()
	InitMockShard()

	_, err := CreateShardByName(MockCtx, "test_shard")

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestGetAllShards(t *testing.T) {
	InitMockCtx()

	_, err := GetAllShards(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestClusterGetShards(t *testing.T) {
	InitMockCtx()

	_, err := ClusterGetShards(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}
