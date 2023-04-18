package resource

import "github.com/databonfire/bonfire/filter"

type (
	ListRequest struct {
		FilterJsonlized string        `json:"f"`
		Filter          filter.Filter `json:"filter"`
		PerPage         int64         `json:"per_page"`
		Paged           int64         `json:"paged"`
		//Sorts           []*Sort       `json:"sorts"`
		Sort  string `json:"sort"`
		Order string `json:"order"`
	}

	Sort struct {
		By    string `json:"by"`
		Order string `json:"order"`
	}
)

type (
	ListResponse struct {
		Data       []interface{} `json:"data"`
		Pagination *Pagination   `json:"pagination"`
	}

	Pagination struct {
		Total   int64 `json:"total"`
		PerPage int64 `json:"per_page"`
		Paged   int64 `json:"paged"`
	}
)
