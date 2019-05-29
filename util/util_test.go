// +build unit

package util

import (
	"testing"
	"strings"
)

//
func TestUniqueIntArray(t *testing.T) {

	unq := NewUniqueIntArray()

	unq.Append(1)
	unq.Append(1)

	if len(unq.Data) > 1 {
		t.Error("Failed to keep array unique")
	}

	unq.Append(2)
	unq.Append(2)
	unq.Append(3)
	unq.Append(4)

	if len(unq.Data) > 4 {
		t.Error("Failed to keep array unique. Array len should be 4")
	}
}

//
func TestWhereAmI(t *testing.T) {
	info := WhereAmI()
	parts := strings.Split(info, "~")

	if parts[0] != "util_test" {
		t.Errorf("WhereAmI doesn't know where I am. Expected: 'util_test' Got: '%s'", parts[0])
	}
}
