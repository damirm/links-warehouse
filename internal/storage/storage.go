package storage

import (
	"context"

	"github.com/damirm/links-warehouse/internal/model"
)

type TranscationManager interface {
	Transaction(ctx context.Context, fn func(context.Context, Storage) error) error
}

type Storage interface {
	SaveLink(context.Context, model.Link) error
}
