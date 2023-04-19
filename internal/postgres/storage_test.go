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
