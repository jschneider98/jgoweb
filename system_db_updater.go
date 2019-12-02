package jgoweb

import(
	"sync"
	"errors"
	jgowebDb "github.com/jschneider98/jgoweb/db"
	"github.com/jschneider98/jgoweb/util"
	"github.com/gocraft/dbr"
)

// 
type SystemDbUpdater struct {
	Db *jgowebDb.Collection
	DbUpdates []SystemDbUpdateInterface
	DryRun bool
}

// 
func NewSystemDbUpdater(db *jgowebDb.Collection, updates []SystemDbUpdateInterface, dryRun bool) *SystemDbUpdater {
	sdu := &SystemDbUpdater{db, updates, dryRun}

	return sdu
}

//
func (sdu *SystemDbUpdater) SetDebug(debug bool) {
	util.Debug = debug
}

// Get update info for all DBs
func (sdu *SystemDbUpdater) GetDbUpdateInfo() (map[string][]SystemDbUpdateInterface, error) {
	info := make(map[string][]SystemDbUpdateInterface)

	for dbName, dbConn := range sdu.Db.GetConns() {
		ctx := NewContext(sdu.Db)
		ctx.SetDbSession(dbConn.NewSession(nil))
		
		info[dbName] = make([]SystemDbUpdateInterface, 0)

		for _, update := range sdu.DbUpdates {
			up, err := CreateSystemDbUpdateByUpdateName(ctx, update.GetUpdateName())

			if err != nil {
				return nil, err
			}

			if !up.Description.Valid {
				up.SetDescription(update.GetDescription())
			}

			info[dbName] = append(info[dbName], up)
		}
	}

	return info, nil
}

//
func (sdu *SystemDbUpdater) RunAll() error {
	util.Debugln("Starting DB updater...")

	var errcList []<-chan error

	for dbName, dbConn := range sdu.Db.GetConns() {
		util.Debugln("Applying updates for " + dbName)

		errc := sdu.RunAllByDbSession(dbConn.NewSession(nil), dbName)
		errcList = append(errcList, errc)
	}

	return sdu.WaitForPipeline(errcList...)
}

//
func (sdu *SystemDbUpdater) RunAllByDbSession(dbSess *dbr.Session, dbName string) (<-chan error) {
	var err error
	// var tx *dbr.Tx
	errc := make(chan error, 1)

	ctx := NewContext(sdu.Db)
	ctx.SetDbSession(dbSess)

	if sdu.DryRun {
		_, err = ctx.Begin()

		if err != nil {
			errc <- err
			defer close(errc)
			
			return errc
		}
	}

	go func() {
		defer close(errc)

		for _, update := range sdu.DbUpdates {
			// Must clone/copy original update for goroutine to work
			up := update.Clone()
			up.SetContext(ctx)

			util.Debugf("Applying %s: '%s'\n", dbName, update.GetUpdateName())

			err = sdu.Run(up, dbName)

			if err != nil {
				err = errors.New("ERROR: " + dbName + ": '" + update.GetUpdateName() + "': " + err.Error())

				errc <- err
				return
			}
		}

		if sdu.DryRun {
			util.Debugln(dbName + ": Dry Run. Rolling back changes.")

			err = ctx.Rollback()

			if err != nil {
				errc <- err
				return
			}
		}
	}()

	return errc
}

//
func (sdu *SystemDbUpdater) Run(update SystemDbUpdateInterface, dbName string) error {
	needsToRun, err := update.NeedsToRun()

	if err != nil {
		return err
	}

	if needsToRun {
		err := update.Run()

		if err != nil {
			return err
		}

		err = update.SetComplete()

		if err != nil {
			return err
		}

		util.Debugf("%s: '%s' done.\n", dbName, update.GetUpdateName())
	} else {
		util.Debugf("%s: '%s' already applied. Skipping.\n", dbName, update.GetUpdateName())
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
