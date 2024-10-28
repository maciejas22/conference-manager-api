package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"time"

	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/internal/services"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *queryResolver) News(ctx context.Context, page models.Page) (*models.NewsPage, error) {
	news, meta, err := services.GetNews(ctx, r.dbClient, &page)
	if err != nil {
		return nil, err
	}

	return &models.NewsPage{
		Data: news,
		Meta: meta,
	}, nil
}

func (r *queryResolver) TermsAndConditions(ctx context.Context) (*models.TermsOfService, error) {
	tos, err := services.GetTermsAndConditions(ctx, r.dbClient)
	if err != nil {
		return nil, gqlerror.Errorf("Error getting terms of service")
	}

	updatedAt, _ := time.Parse(time.RFC3339, tos.UpdatedAt)
	newSections := make([]*models.Section, len(tos.Sections))
	for i, originalSection := range tos.Sections {
		newSubsections := make([]*models.SubSection, len(originalSection.Subsections))
		for j, originalSubsection := range originalSection.Subsections {
			content := ""
			if originalSubsection.Content != nil {
				content = *originalSubsection.Content
			}
			newSubsections[j] = &models.SubSection{
				ID:      originalSubsection.Id,
				Title:   originalSubsection.Title,
				Content: content,
			}
		}

		newSections[i] = &models.Section{
			ID:          originalSection.Id,
			Title:       &originalSection.Title,
			Content:     originalSection.Content,
			Subsections: newSubsections,
		}
	}

	return &models.TermsOfService{
		ID:              tos.Id,
		UpdatedAt:       updatedAt,
		Introduction:    tos.Introduction,
		Acknowledgement: tos.Acknowledgement,
		Sections:        newSections,
	}, nil
}
