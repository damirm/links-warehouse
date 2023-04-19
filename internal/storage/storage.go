package storage

import (
	"context"

	"github.com/damirm/links-warehouse/internal/model"
)

type Storage interface {
	SaveLink(context.Context, model.Link) error
}
