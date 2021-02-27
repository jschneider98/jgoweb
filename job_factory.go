package jgoweb

import (
	"errors"
	"fmt"
	"net/url"
)

type JobFactoryInterface interface {
	New(ContextInterface, string, url.Values) (JobInterface, error)
}

//
type JobFactoryExample struct{}

//
func (jf *JobFactoryExample) New(ctx ContextInterface, name string, params url.Values) (JobInterface, error) {
	switch name {
	case "test":
		// NOTE: ctx/params not needed for job example, but may be needed for other jobs
		return NewJobExample(), nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid job: %s", name))
	}
}
