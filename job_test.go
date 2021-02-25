// +build unit

package jgoweb

import (
	"testing"
	"time"
)

//
func TestJobExample(t *testing.T) {
	j := NewJobExample()
	j.Run()

	time.Sleep(200 * time.Millisecond)

	if j.NumSleeps < 2 {
		t.Errorf("Number of sleeps is less than 2 (%v)", j.NumSleeps)
	}
}

//
func TestJobExampleQuit(t *testing.T) {
	j := NewJobExample()
	j.Run()
	j.Quit()

	time.Sleep(200 * time.Millisecond)

	if j.NumSleeps > 1 {
		t.Errorf("Number of sleeps is greater than 1 (%v)", j.NumSleeps)
	}
}
