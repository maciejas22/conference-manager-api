package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/maciejas22/conference-manager/api/db"
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

type NewsRepository interface {
	GetAll() ([]News, error)
}

type newsRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewNewsRepository(ctx context.Context, db *db.DB) NewsRepository {
	return &newsRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *newsRepository) GetAll() ([]News, error) {
	var news []News
	n := &News{}
	query := "SELECT id, title, content, created_at FROM " + n.TableName() + " ORDER BY created_at DESC"
	err := r.db.SqlConn.Select(
		&news,
		query,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("could not get news")
	}
	return news, nil
}
