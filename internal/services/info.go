package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converter"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetNews(ctx context.Context, db *db.DB) ([]*models.News, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	news, err := repositories.GetAllNews(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result []*models.News
	for _, n := range news {
		result = append(result, converter.ConvertNewsRepoToSchema(&n))
	}

	return result, nil
}

func GetTermsAndConditions(ctx context.Context, db *db.DB) (*models.TermsOfService, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	termsOfService, err := repositories.GetTermsOfService(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converter.ConvertTosRepoToSchema(&termsOfService), nil
}

func GetToSSections(ctx context.Context, db *db.DB, tosId string) ([]*models.Section, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	sections, err := repositories.GetToSSections(tx, tosId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result []*models.Section
	for _, s := range sections {
		result = append(result, converter.ConvertSectionRepoToSchema(&s))
	}

	return result, nil
}

func GetToSSubsections(ctx context.Context, db *db.DB, sectionId string) ([]*models.SubSection, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	subSections, err := repositories.GetToSSubsections(tx, sectionId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result []*models.SubSection
	for _, s := range subSections {
		result = append(result, converter.ConvertSubsectionRepoToSchema(&s))
	}

	return result, nil
}
