package jgoweb

import (
	"github.com/carlescere/scheduler"
	"github.com/jschneider98/jgoweb/util"
	"log"
)

// SystemJob = Data store backed job (job status, queued, started, ended etc)
// Job = Actual job to run
// SchedJob = scheduler job that checks SystemJobs and manages running Jobs
// @TODO: Add cancel Job functionality (like 80% of the way there already with Job.Quit())

type JobQueue struct {
	MaxConcurrency  uint64
	ProcessInterval int
	SchedJob        *scheduler.Job
	Debug           bool
	Ctx             ContextInterface
	jobs            []JobInterface
	dataStore       JobQueueStoreInterface
	factory         JobFactoryInterface
}

//
func NewJobQueue(ctx ContextInterface, dataStore JobQueueStoreInterface, factory JobFactoryInterface) (*JobQueue, error) {
	jq := &JobQueue{Ctx: ctx, dataStore: dataStore, factory: factory}
	jq.MaxConcurrency = 50
	jq.jobs = make([]JobInterface, 0)

	// Num seconds to process jobs
	jq.ProcessInterval = 60

	return jq, nil
}

//
func (jq *JobQueue) Run() error {
	var err error

	if jq.SchedJob != nil {
		jq.Stop()
	}

	fn := func() {
		err := jq.ProcessJobs()

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
		}
	}

	jq.SchedJob, err = scheduler.Every(jq.ProcessInterval).Seconds().Run(fn)

	return err
}

// Same as normal, but worker will process only one job at a time
func (jq *JobQueue) RunWorker() error {
	var err error

	if jq.SchedJob != nil {
		jq.Stop()
	}

	fn := func() {
		err := jq.WorkerProcessJobs()

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
		}
	}

	jq.SchedJob, err = scheduler.Every(jq.ProcessInterval).Seconds().Run(fn)

	return err
}

//
func (jq *JobQueue) Stop() {
	if jq.SchedJob == nil {
		return
	}

	jq.SchedJob.Quit <- true
	jq.SchedJob = nil
}

// Simple wrapper
func (jq *JobQueue) EnqueueJob(job *SystemJob) error {
	return jq.dataStore.EnqueueJob(job)
}

//
func (jq *JobQueue) WorkerProcessJobs() error {
	sysJobs, err := jq.dataStore.GetNextJobs(jq.MaxConcurrency)

	if err != nil {
		return err
	}

	if jq.Debug {
		log.Printf("DEBUG:\n%s\nNum jobs to run: %v\n", util.WhereAmI(), len(sysJobs))
	}

	if sysJobs != nil && len(sysJobs) > 0 {
		// NOTE: Must use distinct DB session per job
		sysJobs[0].Ctx = jq.NewContext()
		jq.processJob(&sysJobs[0])
	}

	return nil
}

//
func (jq *JobQueue) ProcessJobs() error {
	sysJobs, err := jq.dataStore.GetNextJobs(jq.MaxConcurrency)

	if err != nil {
		return err
	}

	if jq.Debug {
		log.Printf("DEBUG:\n%s\nNum jobs to run: %v\n", util.WhereAmI(), len(sysJobs))
	}

	for _, sysJob := range sysJobs {
		// NOTE: Must use distinct DB session per job
		sysJob.Ctx = jq.NewContext()
		jq.processJob(&sysJob)
	}

	return nil
}

//
func (jq *JobQueue) NewContext() ContextInterface {
	ctx := NewContext(jq.Ctx.GetDb())
	ctx.SetDbSession(jq.Ctx.GetDbSession().Connection.NewSession(nil))

	return ctx
}

//
func (jq *JobQueue) processJob(sysJob *SystemJob) error {
	params, err := sysJob.GetDataValues()

	if err != nil {
		err = sysJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return nil
	}

	ctx := jq.NewContext()

	job, err := jq.factory.New(ctx, sysJob.GetName(), params)

	if err != nil {
		err = sysJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return nil
	}

	err = sysJob.Start()

	if err != nil {
		err = sysJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return nil
	}

	go func(job JobInterface, sysJob *SystemJob, debug bool) {
		var err error
		sysJob.Ctx.SetDbSession(sysJob.Ctx.GetDbSession().Connection.NewSession(nil))

		err = job.Run()

		if err != nil {
			err = sysJob.Fail(err)

			if err != nil {
				log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			}

			return
		}

		for {
			select {
			case <-job.GetDoneChannel():
				err = job.GetError()

				if debug {
					log.Printf("DEBUG:\n%s\n%s done.", util.WhereAmI(), sysJob.GetDescription())
				}

				if err != nil {
					err = sysJob.Fail(err)

					if err != nil {
						log.Printf("ERROR: %s %s", util.WhereAmI(), err)
					}
				} else {
					err = sysJob.End()

					if err != nil {
						log.Printf("ERROR: %s %s", util.WhereAmI(), err)
					}
				}

				return
			case <-job.GetCheckinChannel():
				err = sysJob.Checkin(job.GetStatus())

				if err != nil {
					log.Printf("ERROR: %s %s", util.WhereAmI(), err)
				}

			default:
			}
		}
	}(job, sysJob, jq.Debug)

	return nil
}
