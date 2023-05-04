package warehouse

import (
	"context"
	"net/url"
)

type Storage interface {
	SaveLink(context.Context, *Link) error
	LinkExists(context.Context, *url.URL) (bool, error)
	EnqueueURL(context.Context, *url.URL) error
	DequeueURL(context.Context) (*url.URL, error)
	DeleteProcessedURL(context.Context, *url.URL) error

	Transaction(ctx context.Context, fn func(context.Context, Storage) error) error
}
