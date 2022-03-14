// +build unit

package util

import (
	"strings"
	"testing"
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

//
func TestHtmlTemplateToString(t *testing.T) {
	type Inventory struct {
		Material string
		Count    uint
	}

	sweaters := Inventory{"wool", 17}
	tmpl := "{{{.Count}}} items are made of {{{.Material}}}"

	result, err := HtmlTemplateToString(tmpl, sweaters)

	if err != nil {
		t.Error(err)
	}

	if result != "17 items are made of wool" {
		t.Errorf("Incorrect string: %s", result)
	}
}
