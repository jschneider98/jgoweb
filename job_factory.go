package jgoweb

import (
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
	default:
		// NOTE: ctx/params not needed for job example, but may be needed for other jobs
		return NewJobExample(), nil
	}
}
