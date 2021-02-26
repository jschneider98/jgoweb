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
	var err error
	results := make([]SystemJob, 0)

	// force a default max
	if maxConcurrency == 0 {
		maxConcurrency = 100
	}

	if jqs.GetRunningJobs() >= maxConcurrency {
		return results, nil
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

//
func (jqs *JobQueueNativeStore) GetRunningJobs() uint64 {
	var count uint64

	query := `
		SELECT count(*) running_jobs
		FROM system.jobs
		WHERE started_at IS NOT NULL
			AND ended_at IS NULL
		LIMIT 1
	`

	stmt, err := jqs.Ctx.Prepare(query)

	if err != nil {
		return count
	}

	defer stmt.Close()

	stmt.QueryRow().Scan(&count)

	return count
}
