package services

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	filters "github.com/maciejas22/conference-manager/api/internal/db/repositories/shared"
	ToSRepo "github.com/maciejas22/conference-manager/api/internal/db/repositories/terms_of_service"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetNews(ctx context.Context, dbClient *db.DB, p *models.Page) ([]*models.News, *models.NewsMeta, error) {
	var news []repositories.News
	var meta filters.PaginationMeta
	var page filters.Page
	if p == nil {
		page = filters.Page{
			PageNumber: 1,
			PageSize:   10,
		}
	} else {
		page = filters.Page{
			PageNumber: p.Number,
			PageSize:   p.Size,
		}
	}

	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		news, meta, err = repositories.GetNews(tx, page)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	var result []*models.News
	for _, n := range news {
		parsedDate, err := time.Parse(time.RFC3339, n.CreatedAt)
		if err != nil {
			return nil, nil, err
		}

		result = append(result,
			&models.News{
				ID:      n.Id,
				Title:   n.Title,
				Content: n.Content,
				Date:    parsedDate,
			})
	}

	return result, &models.NewsMeta{
		Page: &models.PageInfo{
			TotalItems: meta.TotalItems,
			TotalPages: meta.TotalPages,
			Number:     meta.PageNumber,
			Size:       meta.PageSize,
		},
	}, nil
}

func GetTermsAndConditions(ctx context.Context, dbClient *db.DB) (*ToSRepo.TermsOfService, error) {
	var termsOfService *ToSRepo.TermsOfService
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		termsOfService, err = ToSRepo.GetTermsOfService(tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return &ToSRepo.TermsOfService{}, err
	}

	return termsOfService, nil
}
