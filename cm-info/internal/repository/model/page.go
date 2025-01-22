package filter

type Page struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

type PaginationMeta struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
