package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetNews(ctx context.Context, dbClient *db.DB) ([]*models.News, error) {
	news, err := repositories.GetAllNews(dbClient.QueryExecutor)
	if err != nil {
		return nil, err
	}

	var result []*models.News
	for _, n := range news {
		result = append(result, converters.ConvertNewsRepoToSchema(&n))
	}

	return result, nil
}

func GetTermsAndConditions(ctx context.Context, dbClient *db.DB) (*models.TermsOfService, error) {
	termsOfService, err := repositories.GetTermsOfService(dbClient.QueryExecutor)
	if err != nil {
		return nil, err
	}

	return converters.ConvertTosRepoToSchema(&termsOfService), nil
}

func GetToSSections(ctx context.Context, dbClient *db.DB, tosId int) ([]*models.Section, error) {
	sections, err := repositories.GetToSSections(dbClient.QueryExecutor, tosId)
	if err != nil {
		return nil, err
	}

	var result []*models.Section
	for _, s := range sections {
		result = append(result, converters.ConvertSectionRepoToSchema(&s))
	}

	return result, nil
}

func GetToSSubsections(ctx context.Context, dbClient *db.DB, sectionId int) ([]*models.SubSection, error) {
	subSections, err := repositories.GetToSSubsections(dbClient.QueryExecutor, sectionId)
	if err != nil {
		return nil, err
	}

	var result []*models.SubSection
	for _, s := range subSections {
		result = append(result, converters.ConvertSubsectionRepoToSchema(&s))
	}

	return result, nil
}
