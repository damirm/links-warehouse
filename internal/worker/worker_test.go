package worker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/damirm/links-warehouse/internal/worker"
)

func TestRunJobs(t *testing.T) {
	workersCount := uint(5)
	conf := &worker.Config{
		WorkerCount:   workersCount,
		JobBufferSize: 0,
	}
	worker := worker.NewWorker(conf)

	wg := sync.WaitGroup{}
	jobsCount := workersCount * 2
	wg.Add(int(jobsCount))
	for i := 0; i < int(jobsCount); i++ {
		go worker.Run(func() error {
			wg.Done()
			return nil
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	go worker.StartAndJoin()

	select {
	case <-time.After(5 * time.Second):
		t.Error("job has not been called within 5 seconds")
	case <-done:
	}
}
