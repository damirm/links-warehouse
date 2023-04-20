package processor_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/damirm/links-warehouse/internal/fetcher"
	"github.com/damirm/links-warehouse/internal/parser"
	"github.com/damirm/links-warehouse/internal/processor"
	"github.com/damirm/links-warehouse/internal/storage"
	"github.com/damirm/links-warehouse/internal/worker"
	"github.com/stretchr/testify/require"
)

func TestLinkProcesor(t *testing.T) {
	s := storage.NewInMemoryStorage()
	w := worker.NewWorker(&worker.Config{WorkerCount: 10})
	// TODO: Better to use mocked fetcher.
	f := &fetcher.HttpFetcher{}
	p := &parser.HabrParser{}
	lp := processor.NewLinkProcessor(s, w, f, p)

	w.Start()
	lp.Start()

	go func() {
		<-time.After(5 * time.Second)
		lp.Stop()
	}()

	ctx := context.Background()
	u, err := url.Parse("https://habr.com/ru/companies/skillfactory/articles/729924/")
	require.NoError(t, err)
	queueSize := 10
	for i := 0; i < queueSize; i++ {
		err = s.EnqueueURL(ctx, u)
		require.NoError(t, err)
	}

	// Wait when link processor will be stopped.
	lp.Join()

	w.Stop()
	w.Join()

	if uint64(queueSize) != lp.ProcessedJobs() {
		t.Errorf("failed to process all jobs, expected: %d but got %d", queueSize, lp.ProcessedJobs())
	}
	t.Logf("processed jobs: %d", lp.ProcessedJobs())
}
