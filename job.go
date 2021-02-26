package jgoweb

import (
	"fmt"
	"time"
)

type JobInterface interface {
	Run() error
	Quit()
	IsDone() bool
	GetError() error
	GetStatus() string
	GetDoneChannel() chan bool
	GetCheckinChannel() chan bool
}

type JobExample struct {
	NumSleeps int
	quit      chan bool
	Done      chan bool
	Checkin   chan bool
	status    string
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
func (j *JobExample) Run() error {

	if j.isRunning || j.isDone {
		return nil
	}

	j.quit = make(chan bool, 1)
	j.Done = make(chan bool, 1)
	j.Checkin = make(chan bool, 1)
	j.isRunning = true

	go func(j *JobExample) {
		fmt.Println("1st sleep")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("1st sleep done")

		j.checkin("50% complete")
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
		j.isRunning = false
		j.isDone = true
		j.Done <- true
	}(j)

	return nil
}

//
func (j *JobExample) Quit() {
	j.isRunning = false
	j.isDone = true
	j.Done <- true
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

//
func (j *JobExample) GetStatus() string {
	return j.status
}

//
func (j *JobExample) GetDoneChannel() chan bool {
	return j.Done
}

//
func (j *JobExample) GetCheckinChannel() chan bool {
	return j.Checkin
}

//
func (j *JobExample) checkin(status string) {
	j.status = status
	j.Checkin <- true
}
