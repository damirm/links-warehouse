package parser

import (
	"net/url"

	"github.com/damirm/links-warehouse/internal/model"
)

type Parser interface {
	Parse(*url.URL, string) (*model.Link, error)
}
