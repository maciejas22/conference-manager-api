package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	table "github.com/maciejas22/conference-manager-api/cm-info/internal/repository"
	filter "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/model"
)

type NewsRepo struct {
	Db *sqlx.DB
}

func NewNewsRepo(db *sqlx.DB) NewsRepoInterface {
	return &NewsRepo{Db: db}
}

func (r *NewsRepo) GetNews(page filter.Page) ([]News, filter.PaginationMeta, error) {
	var news []News

	offset := (page.PageNumber - 1) * page.PageSize

	var totalItems int
	countQuery := "SELECT COUNT(*) FROM " + table.GetTableName(table.News)
	err := r.Db.Get(&totalItems, countQuery)
	if err != nil {
		return nil, filter.PaginationMeta{}, errors.New("Could not get total items")
	}

	totalPages := (totalItems + page.PageSize - 1) / page.PageSize

	query := `
		SELECT id, title, content, created_at
		FROM ` + table.GetTableName(table.News) + `
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2 
	`
	err = r.Db.Select(&news, query, page.PageSize, offset)
	if err != nil {
		return nil, filter.PaginationMeta{}, errors.New("Could not get news")
	}

	paginationMeta := filter.PaginationMeta{
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return news, paginationMeta, nil
}
