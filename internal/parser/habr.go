package parser

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/damirm/links-warehouse/internal/model"
)

func Parse(u *url.URL, body string) (*model.Link, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	title := doc.Find(".tm-title span").Text()

	commentsCount, err := strconv.Atoi(strings.TrimSpace(doc.Find(".tm-article-comments-counter-link__value").First().Text()))
	if err != nil {
		// TODO: Skip for now.
	}

	return &model.Link{
		URL:           u,
		Title:         title,
		CommentsCount: int64(commentsCount),
	}, nil
}
