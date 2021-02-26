// +build integration

package jgoweb

import (
	"testing"
)

//
func TestJobQueueNativeStoreGetNextJobs(t *testing.T) {
	InitMockCtx()
	jqs, err := NewJobQueueNativeStore(MockCtx)

	if err != nil {
		t.Errorf("ERROR: (%v)", err)
	}

	_, err = jqs.GetNextJobs(100)

	if err != nil {
		t.Errorf("ERROR: (%v)", err)
	}
}
