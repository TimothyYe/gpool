package gpool

import "fmt"

// Job represents user request, function which should be executed in some worker.
type Job func()

// Gorouting instance which can accept client jobs
type worker struct {
	workerPool chan *worker
	jobChannel chan Job
	stop       chan bool
}

// Try calls function and handles the panic.
func Try(fun func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recover from panic, err is: %s \r\n", err)
		}
	}()
	fun()
}

func (w *worker) start() {
	go func() {
		var job Job
		for {
			// add worker to the pool
			w.workerPool <- w

			select {
			case job = <-w.jobChannel:
				// once job is executed, worker will be put back to the poll again
				Try(job)
			case <-w.stop:
				w.stop <- true
				return
			}
		}
	}()
}

func newWorker(pool chan *worker) *worker {
	return &worker{
		workerPool: pool,
		jobChannel: make(chan Job),
		stop:       make(chan bool),
	}
}
