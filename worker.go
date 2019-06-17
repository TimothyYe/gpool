package gpool

// Job represents user request, function which should be executed in some worker.
type Job func()

// Gorouting instance which can accept client jobs
type worker struct {
	workerPool chan *worker
	jobChannel chan Job
	stop       chan bool
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
				job()
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
