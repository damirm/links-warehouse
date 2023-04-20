package postgres_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/damirm/links-warehouse/internal/model"
	"github.com/damirm/links-warehouse/internal/postgres"
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

	testLink := model.Link{
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
