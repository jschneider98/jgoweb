// +build integration

package jgoweb

import (
	"testing"
)

//
func TestNewSearchParams(t *testing.T) {
	sp := NewSearchParams()

	_, err := sp.BuildDefaultCondition()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	sp.Query = "test"

	_, err = sp.BuildDefaultCondition()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	sp.Query = "test two"

	_, err = sp.BuildDefaultCondition()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}
}
