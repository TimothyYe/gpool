package gpool

import "sync"

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
