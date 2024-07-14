package models

import (
	filters "github.com/maciejas22/conference-manager/api/db/repositories/shared"
)

type Page struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}

func (p *Page) ToRepo() filters.Page {
	if p == nil {
		return filters.Page{
			PageNumber: 1,
			PageSize:   10,
		}
	}

	return filters.Page{
		PageNumber: p.Number,
		PageSize:   p.Size,
	}
}

type PageInfo struct {
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Number     int `json:"number"`
	Size       int `json:"size"`
}

func (p *PageInfo) ToRepo() *filters.PaginationMeta {
	if p == nil {
		return nil
	}

	return &filters.PaginationMeta{
		PageNumber: p.Number,
		PageSize:   p.Size,
		TotalItems: p.TotalItems,
		TotalPages: p.TotalPages,
	}
}

func (e Order) ToRepo() filters.Order {
	if e == OrderAsc {
		return filters.ASC
	}
	return filters.DESC
}

type Sort struct {
	Column string `json:"column"`
	Order  Order  `json:"order"`
}

func (s *Sort) ToRepo() *filters.Sort {
	if s == nil {
		return nil
	}

	return &filters.Sort{
		Column: s.Column,
		Order:  s.Order.ToRepo(),
	}
}
