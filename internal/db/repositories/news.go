package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	filters "github.com/maciejas22/conference-manager/api/internal/db/repositories/shared"
)

type News struct {
	Id        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func (n *News) TableName() string {
	return "public.news"
}

func GetNews(tx *sqlx.Tx, page filters.Page) ([]News, filters.PaginationMeta, error) {
	var news []News

	offset := (page.PageNumber - 1) * page.PageSize
	log.Println("offset", offset)

	var totalItems int
	countQuery := "SELECT COUNT(*) FROM " + (new(News)).TableName()
	err := tx.Get(&totalItems, countQuery)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	totalPages := (totalItems + page.PageSize - 1) / page.PageSize

	query := `
		SELECT id, title, content, created_at
		FROM ` + (new(News)).TableName() + `
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2 
	`
	err = tx.Select(&news, query, page.PageSize, offset)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	paginationMeta := filters.PaginationMeta{
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return news, paginationMeta, nil
}
