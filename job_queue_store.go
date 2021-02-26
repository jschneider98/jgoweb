package jgoweb

type JobQueueStoreInterface interface {
	GetNextJobs(maxConcurrency uint64) ([]SystemJob, error)
	EnqueueJob(*SystemJob) error
}

type JobQueueNativeStore struct {
	Ctx ContextInterface
}

//
func NewJobQueueNativeStore(ctx ContextInterface) (*JobQueueNativeStore, error) {
	jqs := &JobQueueNativeStore{Ctx: ctx}

	return jqs, nil
}

//
func (jqs *JobQueueNativeStore) GetNextJobs(maxConcurrency uint64) ([]SystemJob, error) {
	var results []SystemJob
	var err error

	// force a default max
	if maxConcurrency == 0 {
		maxConcurrency = 100
	}

	stmt := jqs.Ctx.SelectBySql(`
		SELECT *
			FROM system.jobs
			WHERE
				started_at IS NULL
			ORDER BY EXTRACT(EPOCH FROM now() - queued_at)/60 + priority::numeric DESC
		`).
		Limit(maxConcurrency)

	_, err = stmt.Load(&results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

//
func (jqs *JobQueueNativeStore) EnqueueJob(job *SystemJob) error {
	if job.Ctx == nil {
		job.Ctx = jqs.Ctx
	}

	return job.Save()
}
