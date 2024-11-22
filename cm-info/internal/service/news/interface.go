package service

import (
	filter "github.com/maciejas22/conference-manager-api/cm-info/internal/service/model"
)

type NewsServiceInterface interface {
	GetNews(page filter.Page) ([]News, filter.PaginationMeta, error)
}
