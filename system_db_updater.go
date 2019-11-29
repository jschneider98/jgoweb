package jgoweb

import(
	"sync"
	"errors"
	"github.com/jschneider98/jgoweb/util"
	"github.com/gocraft/dbr"
)

// 
type SystemDbUpdateInterface interface {
	NeedsToRun(*dbr.Tx) bool
	Run(*dbr.Tx) error
	GetUpdateName() string
}

// 
type SystemDbUpdater struct {
	DbUpdates []SystemDbUpdateInterface
	DbConns map[string]*dbr.Connection
	DryRun bool
}

// 
func NewSystemDbUpdater(updates []SystemDbUpdateInterface, dbConns map[string]*dbr.Connection) *SystemDbUpdater {
	sdu := &SystemDbUpdater{updates, dbConns, false}

	return sdu
}

//
func (sdu *SystemDbUpdater) RunAll() error {
	util.Debugln("Starting DB updater...")

	var errcList []<-chan error
	// errc := make(chan error, 1)

	for _, dbConn := range sdu.DbConns {
		dbSess := dbConn.NewSession(nil)

		errc := sdu.RunAllByDbSession(dbSess)
		errcList = append(errcList, errc)
	}

	return sdu.WaitForPipeline(errcList...)
}

//
func (sdu *SystemDbUpdater) RunAllByDbSession(dbSess *dbr.Session) (<-chan error) {
	var err error
	var tx *dbr.Tx
	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		for _, update := range sdu.DbUpdates {
			tx, err = dbSess.Begin()

			if err != nil {
				errc <- err
				return
			}

			err = sdu.Run(update, tx)

			if err != nil {
				errc <- err
				return
			}

			err = tx.Commit()

			if err != nil {
				errc <- err
				return
			}
		}
	}()

	return errc
}

//
func (sdu *SystemDbUpdater) Run(update SystemDbUpdateInterface, tx *dbr.Tx) error {
	needsToRun, err := update.NeedsToRun(tx)

	if err != nil {
		return err
	}

	if needsToRun {
		return update.Run(tx)
	}

	return nil
}

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
// Based on https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d
func (sdu *SystemDbUpdater) MergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	// We must ensure that the output channel has the capacity to hold as many errors
	// as there are error channels. This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan error) {

		for n := range c {
			out <- n
		}

		wg.Done()
	}

	wg.Add(len(cs))
	
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// WaitForPipeline waits for results from all error channels.
func (sdu *SystemDbUpdater) WaitForPipeline(errs ...<-chan error) error {
	var err error
	var msg string

	errc := sdu.MergeErrors(errs...)

	for err := range errc {

		if err != nil {
			msg += err.Error() + "\n"
		}
	}

	if msg != "" {
		err = errors.New(msg)
	}

	return err
}
