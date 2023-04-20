package parser_test

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/damirm/links-warehouse/internal/model"
	"github.com/damirm/links-warehouse/internal/parser"
	"github.com/stretchr/testify/require"
)

//go:embed test_data/article1.html
var articleBody string

func TestHabrParser(t *testing.T) {
	habrParser := parser.HabrParser{}

	u, _ := url.Parse("https://habr.com/en/companies/habr/articles/435764/")
	link, err := habrParser.Parse(u, articleBody)
	require.NoError(t, err)

	publishedAt, err := time.Parse(time.RFC3339Nano, "2019-01-15T11:15:36.000Z")
	require.NoError(t, err)
	expected := model.Link{
		URL:           u,
		Title:         "Hello world! Or Habr in English, v1.0",
		Tags:          []string{"Habr", "cake", "UFO"},
		CommentsCount: 249,
		PublishedAt:   publishedAt,
	}

	if expected.Title != link.Title {
		t.Fatalf("expected title '%s' but got '%s'", expected.Title, link.Title)
	}
	if expected.CommentsCount != link.CommentsCount {
		t.Fatalf("expected comments count '%d' but got '%d'", expected.CommentsCount, link.CommentsCount)
	}
	if !expected.PublishedAt.Equal(link.PublishedAt) {
		t.Fatalf("expected publised at '%s' but got '%s'", expected.PublishedAt, link.PublishedAt)
	}
	assertStringSliceEquals(t, expected.Tags, link.Tags)
}

func assertStringSliceEquals(t *testing.T, a, b []string) {
	if len(a) != len(b) {
		t.Fatalf("slice len is not equal: %d != %d", len(a), len(b))
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("a[%d] = '%s', but b[%d] = '%s'", i, a[i], i, b[i])
		}
	}
}
