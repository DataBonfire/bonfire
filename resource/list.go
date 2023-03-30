package resource

type (
	ListRequest struct {
		FilterJsonlized string  `json:"f"`
		Filter          Filter  `json:"filter"`
		PerPage         int64   `json:"per_page"`
		Paged           int64   `json:"paged"`
		Sorts           []*Sort `json:"sorts"`
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
