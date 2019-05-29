// +build integration

package jgoweb

import (
	"testing"
)

//
func TestNewQueryIntBuilder(t *testing.T) {
	InitMockCtx()
	rawQuery := "test"
	builder := NewQueryIntBuilder(MockCtx, rawQuery)

	if rawQuery != builder.RawQuery {
		t.Errorf("Constructor fail. Expected '%s' Got: '%s'", rawQuery, builder.RawQuery)
	}
}

//
func TestBuild(t *testing.T) {
	rawQuery := "super (barbosa or schneider) and not news"
	builder := NewQueryIntBuilder(MockCtx, rawQuery)

	_, err := builder.Build()

	if err != nil {
		t.Error(err)
	}
}
