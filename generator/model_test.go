// +build integration

package generator

import (
	"fmt"
	"github.com/jschneider98/jgoweb"
	"testing"
)

//
func TestGenerateTest(t *testing.T) {
	mg, err := NewModelGenerator(jgoweb.MockCtx, "public", "accounts", "s")

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	code := mg.GenerateTest()

	fmt.Printf("\n%s\n", code)
}
