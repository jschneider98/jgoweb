package jgoweb

import (
	"fmt"
	"time"
)

type JobInterface interface {
	Run() error
	Quit()
	IsDone() bool
}

type JobParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type JobExample struct {
	NumSleeps int
	quit      chan bool
	isRunning bool
	isDone    bool
}

//
func NewJobExample() *JobExample {
	j := &JobExample{}

	return j
}

//
func (j *JobExample) Run() {

	if j.isRunning || j.isDone {
		return
	}

	j.quit = make(chan bool, 1)
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
	}(j)
}

//
func (j *JobExample) Quit() {
	j.isDone = true
	j.quit <- true
}

//
func (j *JobExample) IsDone() bool {
	return j.isDone
}
