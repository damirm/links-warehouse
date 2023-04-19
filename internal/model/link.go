package model

import (
	"net/url"
	"time"
)

type Link struct {
	URL           url.URL
	Title         string
	Tags          []string
	CommentsCount int64
	ViewsCount    int64
	Rating        int64
	PublishedAt   time.Time
	Author        Author
	Complexity    Complexity
	Status        LinkStatus
}

type Complexity uint8

const (
	ComplexityEasy Complexity = iota
	ComplexityMedium
	ComplexityHard
)

type LinkStatus uint8

const (
	LinkStatusNew LinkStatus = iota
	LinkStatusProcessed

	// Link has been read by user.
	LinkStatusDone
)

type Author struct {
	Name   string
	Login  string
	Rating float64
	Karma  int64
}
