package parser

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/damirm/links-warehouse/internal/warehouse"
)

type HabrParser struct{}

func (p *HabrParser) Parse(u *url.URL, body string) (*warehouse.Link, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	title := doc.Find(".tm-title span").Text()

	commentsCount, err := strconv.Atoi(strings.TrimSpace(doc.Find(".tm-article-comments-counter-link__value").First().Text()))
	if err != nil {
		log.Printf("failed to parse comments count: %v", err)
	}

	var publishedAt time.Time
	datePublished, exists := doc.Find(".tm-article-datetime-published time").Attr("datetime")
	if exists {
		publishedAt, err = time.Parse(time.RFC3339Nano, datePublished)
		if err != nil {
			log.Printf("failed to parse published at: %v", err)
		}
	}

	var tags []string
	doc.Find(".tm-article-presenter__meta .tm-separated-list").First().Find(".tm-separated-list__item a").Each(func(n int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})

	return &warehouse.Link{
		URL:           u,
		Title:         title,
		Tags:          tags,
		CommentsCount: int64(commentsCount),
		PublishedAt:   publishedAt,
	}, nil
}
