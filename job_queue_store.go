package jgoweb

type JobQueueStoreInterface interface {
	GetNextJobs(uint64) ([]QueueJob, error)
	EnqueueJob(*QueueJob) error
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
func (jqs *JobQueueNativeStore) GetNextJobs(maxConcurrency uint64) ([]QueueJob, error) {
	var err error
	results := make([]QueueJob, 0)

	// force a default max
	if maxConcurrency == 0 {
		maxConcurrency = 100
	}

	runningJobs := jqs.GetRunningJobs()

	if runningJobs >= maxConcurrency {
		return results, nil
	}

	limit := maxConcurrency - runningJobs

	stmt := jqs.Ctx.Select("*").
		From("system.jobs").
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
		return 10000000
	}

	defer stmt.Close()

	stmt.QueryRow().Scan(&count)

	return count
}
