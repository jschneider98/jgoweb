// +build unit

package generator

import (
	"fmt"
	"testing"
)

//
func TestReportViewGenerate(t *testing.T) {
	fields := make([]string, 4)

	fields[0] = "id"
	fields[1] = "first_name"
	fields[2] = "last_name"
	fields[3] = "c.test"

	r := NewReportViewGenerator("TestReport", fields)

	fmt.Println(r.GenerateView())
}
