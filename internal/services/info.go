package services

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetNews(ctx context.Context, dbClient *db.DB) ([]*models.News, error) {
	var news []repositories.News
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		news, err = repositories.GetAllNews(tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var result []*models.News
	for _, n := range news {
		parsedDate, err := time.Parse(time.RFC3339, n.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result,
			&models.News{
				ID:      n.Id,
				Title:   n.Title,
				Content: n.Content,
				Date:    parsedDate,
			})
	}

	return result, nil
}

func GetTermsAndConditions(ctx context.Context, dbClient *db.DB) (*models.TermsOfService, error) {
	var termsOfService repositories.TermsOfService
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		termsOfService, err = repositories.GetTermsOfService(tx)
		if err != nil {
			return err
		}

		return nil
	})

	updatedAt, err := time.Parse(time.RFC3339, termsOfService.UpdatedAt)
	if err != nil {
		return &models.TermsOfService{}, err
	}

	return &models.TermsOfService{
		ID:              termsOfService.Id,
		UpdatedAt:       updatedAt,
		Introduction:    termsOfService.Introduction,
		Acknowledgement: termsOfService.Acknowledgement,
	}, err
}

func GetToSSections(ctx context.Context, dbClient *db.DB, tosId int) ([]*models.Section, error) {
	var sections []repositories.Section
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		sections, err = repositories.GetToSSections(tx, tosId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var result []*models.Section
	for _, s := range sections {
		result = append(result, &models.Section{
			ID:      s.Id,
			Title:   &s.Title,
			Content: s.Content,
		})
	}

	return result, nil
}

func GetToSSubsections(ctx context.Context, dbClient *db.DB, sectionId int) ([]*models.SubSection, error) {
	var subsections []repositories.Subsection
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		subsections, err = repositories.GetToSSubsections(tx, sectionId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var result []*models.SubSection
	for _, s := range subsections {
		result = append(result, &models.SubSection{
			ID:      s.Id,
			Title:   s.Title,
			Content: *s.Content,
		})
	}

	return result, nil
}
