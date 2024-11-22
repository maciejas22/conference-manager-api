package repository

import (
	filter "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/model"
)

type NewsRepoInterface interface {
	GetNews(page filter.Page) ([]News, filter.PaginationMeta, error)
}
