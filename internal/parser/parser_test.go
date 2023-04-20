package parser_test

import (
	_ "embed"
	"net/url"
	"testing"

	"github.com/damirm/links-warehouse/internal/model"
	"github.com/damirm/links-warehouse/internal/parser"
	"github.com/stretchr/testify/require"
)

//go:embed test_data/article1.html
var articleBody string

func TestHabrParser(t *testing.T) {
	u, _ := url.Parse("https://habr.com/en/companies/habr/articles/435764/")
	link, err := parser.Parse(u, articleBody)
	require.NoError(t, err)

	expected := model.Link{
		URL:           u,
		Title:         "Hello world! Or Habr in English, v1.0",
		CommentsCount: 249,
	}

	if expected.Title != link.Title {
		t.Fatalf("expected title '%s' but got '%s'", expected.Title, link.Title)
	}
	if expected.CommentsCount != link.CommentsCount {
		t.Fatalf("expected comments count '%d' but got '%d'", expected.CommentsCount, link.CommentsCount)
	}
}
