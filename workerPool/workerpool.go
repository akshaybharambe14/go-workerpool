package workerpool

import (
	"fmt"
	"sync/atomic"
)

// Pool - Worker pool configuration
type Pool struct {
	Jobs       []Job       // Jobs to be executed
	WorkerCnt  int         // Number of workers consuming provided jobs
	ErrCnt     int32       // Number of errors occurred in execution
	JobsChan   chan Job    // Channel to provide jobs to workers
	ResChan    chan Result // Channel to receive result from workers
	IsFinished bool        //
	DebugMode  bool        // if true, print worker info
}

// Job - Single job to be executed
type Job struct {
	ID      int         // job id provided by user, optional
	Execute WorkFunc    // type of function to be executed
	Data    interface{} // data to be provided to function
	Result  Result      // function result
}

// WorkFunc - Worker pool function signature
type WorkFunc func(d interface{}) Result // Worker pool function signature

// Result - Job result
type Result struct {
	Err error
}

// NewPool - Initiate new pool with given configuration
func NewPool(jobs []Job, workers int, debugMode bool) *Pool {
	if jobs == nil {
		jobs = []Job{}
	}
	if workers < 1 {
		workers = 2
	}
	pool := Pool{}
	pool.Jobs = jobs
	pool.WorkerCnt = workers
	pool.JobsChan = make(chan Job, len(jobs))
	pool.ResChan = make(chan Result, len(jobs))
	pool.DebugMode = debugMode
	return &pool
}

// AddJobs - Add more jobs to existing jobs
func (p *Pool) AddJobs(jobs ...Job) {
	p.Jobs = append(p.Jobs, jobs...)
}

// Start - Start worker pool execution
func (p *Pool) Start() {
	for w := 0; w < p.WorkerCnt; w++ {
		go p.addWorker(w)
	}
	for j := 0; j < len(p.Jobs); j++ {
		p.JobsChan <- p.Jobs[j]
	}
	close(p.JobsChan)
	for r := 0; r < len(p.Jobs); r++ {
		p.Jobs[r].Result = <-p.ResChan
	}

	p.IsFinished = true
}

func (p *Pool) addWorker(w int) {
	for j := range p.JobsChan {
		if p.DebugMode {
			fmt.Println(">> Worker ", w, " started job ", j.ID)
		}
		r := j.start(w, p.DebugMode)
		if r.Err != nil {
			atomic.AddInt32(&p.ErrCnt, 1)
		}
		p.ResChan <- r
	}
}

func (j Job) start(w int, debugMode bool) Result {
	if debugMode {
		fmt.Println(">> Worker ", w, " started job ", j.ID)
		defer fmt.Println(">> Worker ", w, " finished job ", j.ID)
	}

	return j.Execute(j.Data)

}
