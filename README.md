# gpool

A lightwegit Goroutine worker pool, modified from [grpool](https://github.com/ivpusic/grpool)

## Example

```go
    runtime.GOMAXPROCS(runtime.NumCPU())

    // new a Goroutine pool with 100 workers and 50 job queue
    pool := gpool.NewPool(100, 50)
    defer pool.Release()

    // set 10 jobs that the pool should wait
    pool.WaitCount(10)

    // submit one or more jobs to pool
    for i := 0; i < 10; i++ {
        count := i

        pool.JobQueue <- func() {
            defer pool.JobDone()
            fmt.Printf("Pool job: %d\n", count)
        }
    }

    // wait until we call JobDone for all jobs
    pool.WaitAll()
```
