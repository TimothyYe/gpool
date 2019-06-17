package gpool

import "sync"

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
			// worker free, add it to pool
			w.workerPool <- w

			select {
			case job = <-w.jobChannel:
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

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan bool
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			worker.jobChannel <- job
		case <-d.stop:
			for i := 0; i < cap(d.workerPool); i++ {
				worker := <-d.workerPool

				worker.stop <- true
				<-worker.stop
			}

			d.stop <- true
			return
		}
	}
}

func newDispatcher(workerPool chan *worker, jobQueue chan Job) *dispatcher {
	d := &dispatcher{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		stop:       make(chan bool),
	}

	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool)
		worker.start()
	}

	go d.dispatch()
	return d
}

// Job represents user request, function which should be executed in some worker.
type Job func()

// Pool represents the pool struct
type Pool struct {
	JobQueue   chan Job
	dispatcher *dispatcher
	wg         sync.WaitGroup
}

// NewPool returns object contains JobQueue reference, which you can use to send job to pool.
func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan *worker, numWorkers)

	pool := &Pool{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher(workerPool, jobQueue),
	}

	return pool
}

// JobDone should be called every time when job function is done if you are using WaitAll.
func (p *Pool) JobDone() {
	p.wg.Done()
}

// WaitCount means many jobs we should wait when calling WaitAll, it is using WaitGroup Add/Done/Wait
func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

// WaitAll waits for all jobs to finish.
func (p *Pool) WaitAll() {
	p.wg.Wait()
}

// Release will release resources used by pool
func (p *Pool) Release() {
	p.dispatcher.stop <- true
	<-p.dispatcher.stop
}
