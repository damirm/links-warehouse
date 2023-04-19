package storage

import (
	"context"

	"github.com/damirm/links-warehouse/internal/model"
)

type Migrator interface {
	Migrate(context.Context) error
}

type Storage interface {
	SaveLink(context.Context, model.Link) error
}
