package jgoweb

import (
	"fmt"
	"time"
)

type JobInterface interface {
	Run() chan bool
	Quit()
	IsDone() bool
	GetError() error
}

type JobParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type JobExample struct {
	NumSleeps int
	quit      chan bool
	finished  chan bool
	isRunning bool
	isDone    bool
	err       error
}

//
func NewJobExample() *JobExample {
	j := &JobExample{}

	return j
}

//
func (j *JobExample) Run() chan bool {

	if j.isRunning || j.isDone {
		return nil
	}

	j.quit = make(chan bool, 1)
	j.finished = make(chan bool, 1)
	j.isRunning = true

	go func(j *JobExample) {
		fmt.Println("1st sleep")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("1st sleep done")

		j.NumSleeps++

		select {
		case <-j.quit:
			fmt.Println("Quiting now")
			return
		default:
		}

		fmt.Println("2nd sleep")
		fmt.Println("2nd sleep done")

		j.NumSleeps++
		j.isDone = true
		j.finished <- true
	}(j)

	return j.finished
}

//
func (j *JobExample) Quit() {
	j.isDone = true
	j.finished <- true
	j.quit <- true
}

//
func (j *JobExample) IsDone() bool {
	return j.isDone
}

//
func (j *JobExample) GetError() error {
	return j.err
}
