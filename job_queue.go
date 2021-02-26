package jgoweb

type JobQueue struct {
	MaxConcurrency uint64
	jobs           []JobInterface
	dataStore      JobQueueStoreInterface
	factory        JobFactoryInterface
}

//
func NewJobQueue(dataStore JobQueueStoreInterface, factory JobFactoryInterface) (*JobQueue, error) {
	jq := &JobQueue{dataStore: dataStore, factory: factory}
	jq.MaxConcurrency = 50
	jq.jobs = make([]JobInterface, 0)

	return jq, nil
}

// Simple wrapper
func (jq *JobQueue) EnqueueJob(job *SystemJob) error {
	return jq.dataStore.EnqueueJob(job)
}
