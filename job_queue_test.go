// +build integration

package jgoweb

import (
	"testing"
	"time"
)

//
func TestJobQueueProcessJob(t *testing.T) {
	InitMockCtx()
	jqs, err := NewJobQueueNativeStore(MockCtx)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	fac := &JobFactoryExample{}

	jq, err := NewJobQueue(MockCtx, jqs, fac)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	sysJob, err := NewSystemJob(MockCtx)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	sysJob.SetName("test")
	sysJob.SetDescription("test")

	err = sysJob.Save()

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	err = jq.Run()

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	sysJob, err = FetchSystemJobById(sysJob.Ctx, sysJob.GetId())

	if err != nil {
		t.Errorf("ERROR: %v\n", err)
	}

	if sysJob.GetError() != "" {
		t.Errorf("ERROR: %v\n", sysJob.GetError())
	}

	if sysJob.GetStatus() == "" {
		t.Errorf("ERROR: System Job status is blank, but should be set.\n")
	}

	if sysJob.GetEndedAt() == "" {
		t.Errorf("ERROR: System Job endded at is blank, but should be set.\n")
	}
}
