package postgres_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/damirm/links-warehouse/internal/postgres"
	"github.com/damirm/links-warehouse/internal/warehouse"
	"github.com/stretchr/testify/require"
)

func TestSaveLink(t *testing.T) {
	db, err := postgres.Connect(testDatabaseConfig)
	require.NoError(t, err)

	ctx := context.Background()
	storage, err := postgres.InitStorage(ctx, db, testDatabaseConfig)
	require.NoError(t, err)

	testURL, err := url.Parse("https://google.com")
	require.NoError(t, err)

	testLink := &warehouse.Link{
		URL: testURL,
	}

	err = storage.SaveLink(ctx, testLink)
	require.NoError(t, err)
}

func TestEnqueueDequeueUrl(t *testing.T) {
	db, err := postgres.Connect(testDatabaseConfig)
	require.NoError(t, err)

	ctx := context.Background()
	storage, err := postgres.InitStorage(ctx, db, testDatabaseConfig)
	require.NoError(t, err)

	testURL, err := url.Parse("https://google.com")
	require.NoError(t, err)

	err = storage.EnqueueURL(ctx, testURL)
	require.NoError(t, err)

	nextURL, err := storage.DequeueURL(ctx)
	require.NoError(t, err)

	if testURL.String() != nextURL.String() {
		t.Errorf("expected %s but got %s", testURL, nextURL)
	}
}

func TestTransaction(t *testing.T) {
	db, err := postgres.Connect(testDatabaseConfig)
	require.NoError(t, err)

	ctx := context.Background()
	storage, err := postgres.InitStorage(ctx, db, testDatabaseConfig)
	require.NoError(t, err)

	u, _ := url.Parse("https://google.com")

	storage.Transaction(ctx, func(ctx context.Context, storage warehouse.Storage) error {
		storage.EnqueueURL(ctx, u)
		return fmt.Errorf("something goes wrong")
	})

	u1, err := storage.DequeueURL(ctx)
	require.NoError(t, err)
	require.Nil(t, u1)
}
