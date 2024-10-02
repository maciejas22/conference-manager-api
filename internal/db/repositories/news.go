package repositories

import (
	"github.com/jmoiron/sqlx"
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

func GetAllNews(tx *sqlx.Tx) ([]News, error) {
	var news []News
	n := &News{}
	query := "SELECT id, title, content, created_at FROM " + n.TableName() + " ORDER BY created_at DESC"
	err := tx.Select(&news, query)
	if err != nil {
		return nil, err
	}
	return news, nil
}
