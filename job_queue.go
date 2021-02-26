package jgoweb

type JobQueue struct {
	MaxConcurrency uint64
	Ctx            ContextInterface
	jobs           []JobInterface
	dataStore      JobQueueStoreInterface
	factory        JobFactoryInterface
}

//
func NewJobQueue(ctx ContextInterface, dataStore JobQueueStoreInterface, factory JobFactoryInterface) (*JobQueue, error) {
	jq := &JobQueue{Ctx: ctx, dataStore: dataStore, factory: factory}
	jq.MaxConcurrency = 50
	jq.jobs = make([]JobInterface, 0)

	return jq, nil
}

// Simple wrapper
func (jq *JobQueue) EnqueueJob(job *SystemJob) error {
	return jq.dataStore.EnqueueJob(job)
}

//
func (jq *JobQueue) ProcessJobs() error {
	sysJobs, err := jq.dataStore.GetNextJobs(jq.MaxConcurrency)

	if err != nil {
		return nil
	}

	for _, sysJob := range sysJobs {
		sysJob.Ctx = jq.Ctx
		jq.processJob(&sysJob)
	}

	return nil
}

//
func (jq *JobQueue) processJob(sysJob *SystemJob) error {
	params, err := sysJob.GetDataValues()

	if err != nil {
		err = sysJob.Fail(err)

		// @TODO: Handle DB failure?
		if err != nil {
			return err
		}

		return nil
	}

	job, err := jq.factory.New(jq.Ctx, sysJob.GetName(), params)

	if err != nil {
		err = sysJob.Fail(err)

		// @TODO: Handle DB failure?
		if err != nil {
			return err
		}

		return nil
	}

	go func(job JobInterface, sysJob *SystemJob) {
		var err error
		err = job.Run()

		if err != nil {
			// @TODO: handle DB failure?
			sysJob.Fail(err)
			return
		}

		for {
			select {
			case <-job.GetDoneChannel():
				err = job.GetError()

				if err != nil {
					// @TODO: handle DB failure?
					sysJob.Fail(err)
				} else {
					// @TODO: handle DB failure?
					sysJob.End()
				}

				return
			case <-job.GetCheckinChannel():
				// @TODO: handle DB failure?
				sysJob.Checkin(job.GetStatus())
			default:
			}
		}
	}(job, sysJob)

	return nil
}
