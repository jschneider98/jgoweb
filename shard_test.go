// +build integration

package jgoweb

import (
	"testing"
)

//
func TestFetchShardByAccountId(t *testing.T) {
	InitMockCtx()

	accountId := "00000000-0000-0000-0000-000000000000"

	shard, err := FetchShardByAccountId(MockCtx, accountId)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if shard != nil {
		t.Errorf("\nInvalid account ID should return nil.\n")
	}

	// @TEMP, @TODO: Build a MockAccount. Hard coding accountId for now.
	accountId = "d1d44049-d0b3-4cf3-90ec-b980bf9d1705"

	shard, err = FetchShardByAccountId(MockCtx, accountId)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if shard == nil {
		t.Errorf("\nValid account ID should return a shard.\n")
	}
}
