package repositories

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type News struct {
	Id        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (n *News) TableName() string {
	return "public.news"
}

func GetAllNews(tx *sqlx.Tx) ([]News, error) {
	var news []News
	n := &News{}
	query := "SELECT id, title, content, created_at FROM " + n.TableName() + " ORDER BY created_at DESC"
	err := tx.Select(
		&news,
		query,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("could not get news")
	}
	return news, nil
}
