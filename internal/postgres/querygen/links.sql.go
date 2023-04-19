// Code generated by sqlc. DO NOT EDIT.
// source: links.sql

package querygen

import (
	"context"
	"encoding/json"
	"time"
)

const createLink = `-- name: CreateLink :exec
insert into links
       (url, title, tags, comments_count, views_count, rating, published_at, author, complexity, status)
values ($1,  $2,    $3,   $4,             $5,          $6,     $7,           $8,     $9,         $10)
`

type CreateLinkParams struct {
	Url           string
	Title         string
	Tags          json.RawMessage
	CommentsCount int32
	ViewsCount    int32
	Rating        int32
	PublishedAt   time.Time
	Author        json.RawMessage
	Complexity    int16
	Status        int16
}

func (q *Queries) CreateLink(ctx context.Context, arg CreateLinkParams) error {
	_, err := q.db.ExecContext(ctx, createLink,
		arg.Url,
		arg.Title,
		arg.Tags,
		arg.CommentsCount,
		arg.ViewsCount,
		arg.Rating,
		arg.PublishedAt,
		arg.Author,
		arg.Complexity,
		arg.Status,
	)
	return err
}
