package jgoweb

import (
	"runtime"
)

type JobQueueStoreInterface interface {
	GetNextJobs() ([]QueueJob, error)
	EnqueueJob(*QueueJob) error
}

type JobQueueNativeStore struct {
	Ctx            ContextInterface
	MaxConcurrency uint64
	MaxBatch       uint64
	MaxMem         uint64
}

//
func NewJobQueueNativeStore(ctx ContextInterface) (*JobQueueNativeStore, error) {
	jqs := &JobQueueNativeStore{Ctx: ctx}
	jqs.MaxConcurrency = 100
	jqs.MaxBatch = 10

	// 0 = unlimited, value should be in MB
	jqs.MaxMem = 0

	return jqs, nil
}

//
func (jqs *JobQueueNativeStore) GetNextJobs() ([]QueueJob, error) {
	var err error
	results := make([]QueueJob, 0)

	if jqs.IsMemExceeded() {
		return results, nil
	}

	if jqs.MaxBatch > jqs.MaxConcurrency {
		jqs.MaxBatch = jqs.MaxConcurrency
	}

	runningJobs := jqs.GetRunningJobs()

	if runningJobs >= jqs.MaxConcurrency {
		return results, nil
	}

	limit := jqs.MaxConcurrency - runningJobs

	if limit > jqs.MaxBatch {
		limit = jqs.MaxBatch
	}

	stmt := jqs.Ctx.Select("*").
		From("queue.jobs").
		Where("started_at IS NULL AND ended_at IS NULL").
		OrderBy("EXTRACT(EPOCH FROM now() - queued_at)/60 + priority::numeric DESC").
		Limit(limit)

	_, err = stmt.Load(&results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

//
func (jqs *JobQueueNativeStore) EnqueueJob(job *QueueJob) error {

	if job.Ctx == nil {
		job.Ctx = jqs.Ctx
	}

	return job.Save()
}

//
func (jqs *JobQueueNativeStore) IsMemExceeded() bool {

	if jqs.MaxMem == 0 {
		return false
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if jqs.ByteToMb(m.Alloc) >= jqs.MaxMem {
		return true
	}

	return false
}

//
func (jqs *JobQueueNativeStore) ByteToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//
func (jqs *JobQueueNativeStore) GetRunningJobs() uint64 {
	var count uint64

	query := `
		SELECT count(*) running_jobs
		FROM queue.jobs
		WHERE started_at IS NOT NULL
			AND ended_at IS NULL
		LIMIT 1
	`

	stmt := jqs.Ctx.SelectBySql(query)

	_, err := stmt.Load(&count)

	if err != nil {
		return 10000000
	}

	return count
}
