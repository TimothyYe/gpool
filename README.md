# gpool

A lightwegit Goroutine pool.

## Example

```go
    numCPUs := runtime.NumCPU()
    runtime.GOMAXPROCS(numCPUs)

    // new a Goroutine pool with 100 workers and 50 job queue
    pool := grpool.NewPool(100, 50)
    defer pool.Release()

    // set 10 jobs that the pool should wait
    pool.WaitCount(10)

    // submit one or more jobs to pool
    for i := 0; i < 10; i++ {
        count := i

        pool.JobQueue <- func() {
            defer pool.JobDone()
            fmt.Printf("hello %d\n", count)
        }
    }

    // wait until we call JobDone for all jobs
    pool.WaitAll()
```
