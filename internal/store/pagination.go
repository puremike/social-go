package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PagQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Until  string   `json:"until"`
}

func (fq PagQuery) Parse(r *http.Request) (PagQuery, error) {

	if search := r.URL.Query().Get("search"); search != "" {
		fq.Search = search
	}

	if tags := r.URL.Query().Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ", ")
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		lmt, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = lmt
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		off, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = off
	}


	if sort := r.URL.Query().Get("sort"); sort != "" {
		fq.Sort = sort
	}
	
	return fq, nil
}