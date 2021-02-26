// +build integration

package jgoweb

import (
	"testing"
)

//
func TestNewJobQueue(t *testing.T) {
	InitMockCtx()
	jqs, err := NewJobQueueNativeStore(MockCtx)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	fac := &JobFactoryExample{}

	_, err = NewJobQueue(jqs, fac)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}
