package service

type Page struct {
	PageNumber int
	PageSize   int
}

type PaginationMeta struct {
	PageNumber int
	PageSize   int
	TotalItems int
	TotalPages int
}
