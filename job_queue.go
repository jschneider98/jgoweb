package jgoweb

import (
	"github.com/carlescere/scheduler"
	"github.com/jschneider98/jgoweb/util"
	"log"
)

// QueueJob = Data store backed job (job status, queued, started, ended etc)
// Job = Actual job to run
// SchedJob = scheduler job that checks QueueJobs and manages running Jobs
// @TODO: Add cancel Job functionality (like 80% of the way there already with Job.Quit())

type JobQueue struct {
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
	jq.jobs = make([]JobInterface, 0)

	// Num seconds to process jobs
	jq.ProcessInterval = 5

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

//
func (jq *JobQueue) Stop() {
	if jq.SchedJob == nil {
		return
	}

	jq.SchedJob.Quit <- true
	jq.SchedJob = nil
}

// Simple wrapper
func (jq *JobQueue) EnqueueJob(job *QueueJob) error {
	return jq.dataStore.EnqueueJob(job)
}

//
func (jq *JobQueue) ProcessJobs() error {
	qJobs, err := jq.dataStore.GetNextJobs()

	if err != nil {
		return err
	}

	if jq.Debug {
		log.Printf("DEBUG:\n%s\nNum jobs to run: %v\n", util.WhereAmI(), len(qJobs))
	}

	for _, qJob := range qJobs {
		// NOTE: Use distinct DB session per job
		qJob.Ctx = jq.NewContext()
		go jq.processJob(qJob, jq.Debug)
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
func (jq *JobQueue) processJob(sj QueueJob, debug bool) error {
	qJob := &sj

	if debug {
		log.Printf("DEBUG:\n%s\n%s starting.\n************\n", util.WhereAmI(), qJob.GetDescription())
	}

	params, err := qJob.GetDataValues()

	if err != nil {
		err = qJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return nil
	}

	job, err := jq.factory.New(qJob.Ctx, qJob.GetName(), params)

	if err != nil {
		err = qJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return nil
	}

	err = qJob.Start()

	if err != nil {
		err = qJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return nil
	}

	// qJob.Ctx.SetDbSession(qJob.Ctx.GetDbSession().Connection.NewSession(nil))

	err = job.Run()

	if err != nil {
		err = qJob.Fail(err)

		if err != nil {
			log.Printf("ERROR: %s %s", util.WhereAmI(), err)
			return err
		}

		return err
	}

	for {
		select {
		case <-job.GetCheckinChannel():
			err = qJob.Checkin(job.GetStatus())

			if err != nil {
				log.Printf("ERROR: %s %s", util.WhereAmI(), err)
				return err
			}
		case <-job.GetDoneChannel():
			if debug {
				log.Printf("DEBUG:\n%s\n%s done.\n************\n", util.WhereAmI(), qJob.GetDescription())
			}

			err = job.GetError()

			if err != nil {
				err = qJob.Fail(err)

				if err != nil {
					log.Printf("ERROR: %s %s", util.WhereAmI(), err)
					return err
				}
			} else {
				err = qJob.End()

				if err != nil {
					log.Printf("ERROR: %s %s", util.WhereAmI(), err)
					return err
				}
			}

			return nil
		default:
		}
	}

	return nil
}
