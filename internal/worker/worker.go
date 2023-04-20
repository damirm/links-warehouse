package worker

import (
	"fmt"
	"sync"
)

type Config struct {
	WorkerCount   uint
	JobBufferSize uint
}

type Job func() error

type Worker struct {
	config  *Config
	quit    chan struct{}
	stop    sync.Once
	start   sync.Once
	started bool
	jobs    chan Job
}

func NewWorker(config *Config) *Worker {
	quit := make(chan struct{})
	jobs := make(chan Job, config.JobBufferSize)
	return &Worker{
		config: config,
		quit:   quit,
		jobs:   jobs,
	}
}

func (w *Worker) StartAndJoin() {
	w.start.Do(func() {
		for i := 0; i < int(w.config.WorkerCount); i++ {
			go func(n int) {
				w.work(n)
			}(i)
		}

		w.started = true
	})

	w.Join()
}

func (w *Worker) Run(job Job) {
	w.jobs <- job
}

func (w *Worker) Join() {
	<-w.quit
}

func (w *Worker) Stop() {
	w.stop.Do(func() {
		close(w.quit)
	})
}

// TODO: Use logger instead.
func (w *Worker) work(n int) {
	fmt.Printf("starting worker %d, waiting for tasks...\n", n)

	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				fmt.Printf("jobs channel closed, stopping worker %d...\n", n)
				return
			}
			err := job()
			if err != nil {
				fmt.Printf("job failed: %v\n", err)
			}
		case <-w.quit:
			fmt.Printf("stopping worker %d...\n", n)
			return
		}
	}
}
