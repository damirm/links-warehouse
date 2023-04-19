// Code generated by sqlc. DO NOT EDIT.

package querygen

import (
	"encoding/json"
	"time"
)

type Link struct {
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

type LinksQueue struct {
	Url string
}

type TelegramUser struct {
	ID    int64
	Login string
}