package repositories

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
)

type News struct {
	Id        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func (n *News) TableName() string {
	return "news"
}

func GetAllNews(qe *db.QueryExecutor) ([]News, error) {
	var news []News
	n := &News{}
	query := "SELECT id, title, content, created_at FROM " + n.TableName() + " ORDER BY created_at DESC"
	err := sqlx.Select(
		qe,
		&news,
		query,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("could not get news")
	}
	return news, nil
}
