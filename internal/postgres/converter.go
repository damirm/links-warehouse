package postgres

import (
	"encoding/json"

	"github.com/damirm/links-warehouse/internal/model"
	"github.com/damirm/links-warehouse/internal/postgres/querygen"
)

func insertLinkParams(link model.Link) querygen.InsertLinkParams {
	return querygen.InsertLinkParams{
		Url:           link.URL.String(),
		Title:         link.Title,
		Tags:          toJsonBytes(link.Tags),
		CommentsCount: int32(link.CommentsCount),
		ViewsCount:    int32(link.ViewsCount),
		Rating:        int32(link.Rating),
		PublishedAt:   link.PublishedAt,
		Author:        toJsonBytes(link.Author),
		Complexity:    int16(link.Complexity),
		Status:        int16(link.Status),
	}
}

func toJsonBytes(data any) []byte {
	res, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return res
}
