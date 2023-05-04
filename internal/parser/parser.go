package parser

import (
	"net/url"

	"github.com/damirm/links-warehouse/internal/warehouse"
)

type Parser interface {
	Parse(*url.URL, string) (*warehouse.Link, error)
}
