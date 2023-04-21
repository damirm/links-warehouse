package worker

import (
	"log"
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
	stopped bool
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

func (w *Worker) Start() {
	w.start.Do(func() {
		for i := 0; i < int(w.config.WorkerCount); i++ {
			go func(n int) {
				w.work(n)
			}(i)
		}

		w.started = true
	})
}

func (w *Worker) Run(job Job) {
	if !w.stopped {
		w.jobs <- job
	}
}

func (w *Worker) Join() {
	<-w.quit
}

func (w *Worker) Stop() {
	w.stop.Do(func() {
		w.stopped = true

		close(w.quit)
		close(w.jobs)
	})
}

func (w *Worker) work(n int) {
	log.Printf("starting worker %d, waiting for tasks", n)

	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				log.Printf("jobs channel closed, stopping worker %d", n)
				return
			}
			log.Printf("worker %d picked a job", n)
			err := job()
			if err != nil {
				log.Printf("job failed: %v", err)
			}
			log.Printf("worker %d finished job execution", n)
		case <-w.quit:
			log.Printf("stopping worker %d", n)
			return
		}
	}
}
