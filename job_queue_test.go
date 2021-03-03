// +build integration

package jgoweb

import (
	"testing"
	"time"
)

//
func TestJobQueueProcessJob(t *testing.T) {
	InitMockCtx()
	InitMockUser()
	jqs, err := NewJobQueueNativeStore(MockCtx)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	fac := &JobFactoryExample{}

	jq, err := NewJobQueue(MockCtx, jqs, fac)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	jq.MaxConcurrency = 3
	jq.Debug = true

	// Need a fresh session
	ctx := NewContext(MockCtx.GetDb())
	ctx.SetDbSession(MockCtx.GetDbSession().Connection.NewSession(nil))

	qJob, err := NewQueueJob(ctx)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	qJob.SetAccountId(MockUser.GetAccountId())
	qJob.SetName("test")
	qJob.SetDescription("test")
	qJob.SetPriority("1000000")

	err = qJob.Save()

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	err = jq.Run()

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// qJob, err = FetchQueueJobById(qJob.Ctx, qJob.GetId())

	// if err != nil {
	// 	t.Errorf("ERROR: %v\n", err)
	// }

	// if qJob.GetError() != "" {
	// 	t.Errorf("ERROR: %v\n", qJob.GetError())
	// }

	// if qJob.GetStatus() == "" {
	// 	t.Errorf("ERROR: Queue Job status is blank, but should be set.\n")
	// }

	// if qJob.GetEndedAt() == "" {
	// 	t.Errorf("ERROR: Queue Job endded at is blank, but should be set.\n")
	// }
}
