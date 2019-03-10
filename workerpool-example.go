package main

import (
	"errors"
	"fmt"
	"go-workerPool/workerpool"
	"time"
)

// simple worker pool example
func main() {
	// create jobs
	jobs := make([]workerpool.Job, 100)
	var f workerpool.WorkFunc = func(d interface{}) workerpool.Result {
		r := d.(int)
		if r == 20 || r == 10 || r == 50 {
			return workerpool.Result{
				Err: errors.New("Error"),
			}
		}
		return workerpool.Result{
			Err: nil,
		}
	}
	for i := 0; i < 100; i++ {
		jobs[i] = workerpool.Job{
			ID:      i,
			Data:    i,
			Execute: f,
		}
	}
	// create pool
	p := workerpool.NewPool(jobs, 20, false)
	fmt.Println("Started: ", time.Now())
	// execute worker pool
	p.Start()
	fmt.Println("Exit", time.Now(), " with ", p.ErrCnt, " errors")
}
