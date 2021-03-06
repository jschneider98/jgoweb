// +build unit

package generator

import (
	"fmt"
	"testing"
)

//
func TestReportGenerate(t *testing.T) {
	fields := make([]string, 4)

	fields[0] = "id"
	fields[1] = "first_name"
	fields[2] = "last_name"
	fields[3] = "c.test"

	r := NewReportGenerator("TestReport", fields)

	fmt.Println(r.Generate())
}
