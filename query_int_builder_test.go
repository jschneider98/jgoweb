// +build integration

package jgoweb

import (
	"testing"
)

//
func TestNewQueryIntBuilder(t *testing.T) {
	InitTestCtx()
	rawQuery := "test"
	builder := NewQueryIntBuilder(testCtx, rawQuery)

	if rawQuery != builder.RawQuery {
		t.Errorf("Constructor fail. Expected '%s' Got: '%s'", rawQuery, builder.RawQuery)
	}
}

//
func TestBuild(t *testing.T) {
	rawQuery := "super (barbosa or schneider) and not news"
	builder := NewQueryIntBuilder(testCtx, rawQuery)

	_, err := builder.Build()

	if err != nil {
		t.Error(err)
	}
}
