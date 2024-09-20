package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetParticipantsCount(ctx context.Context, dbClient *db.DB, conferenceId int) (int, error) {
	participantsCount, err := repositories.GetConferenceParticipantsCount(dbClient.QueryExecutor, conferenceId)
	if err != nil {
		return 0, err
	}

	return participantsCount, nil
}

func AddUserToConference(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*int, error) {
	conference, err := repositories.AddConferenceParticipant(dbClient.QueryExecutor, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	cId := int(conference.Id)
	return &cId, nil
}

func RemoveUserFromConference(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*int, error) {
	conference, err := repositories.RemoveConferenceParticipant(dbClient.QueryExecutor, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	cId := int(conference.Id)
	return &cId, nil
}

func IsConferenceParticipant(ctx context.Context, dbClient *db.DB, userId, conferenceID int) (*bool, error) {
	isParticipant, err := repositories.IsConferenceParticipant(dbClient.QueryExecutor, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	return &isParticipant, nil
}

func IsConferenceOrganizer(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*bool, error) {
	isOrganizer, err := repositories.IsConferenceOrganizer(dbClient.QueryExecutor, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	return &isOrganizer, nil
}

func GetOrganizerMetrics(ctx context.Context, dbClient *db.DB, organizerId int) (*models.OrganizerMetrics, error) {
	organizerMetrics, err := repositories.GetOrganizerLevelMetrics(dbClient.QueryExecutor, organizerId)
	if err != nil {
		return nil, err
	}

	return &models.OrganizerMetrics{
		RunningConferences:        organizerMetrics.RunningConferencesCount,
		ParticipantsCount:         organizerMetrics.ParticipantsCount,
		AverageParticipantsCount:  organizerMetrics.AverageParticipantsCount,
		TotalOrganizedConferences: organizerMetrics.TotalOrganizedConferences,
	}, nil
}

func GetParticipantsJoiningTrend(ctx context.Context, dbClient *db.DB, organizerId int) (*models.ParticipantsJoiningTrend, error) {
	participantsJoiningTrend, err := repositories.GetParticipantsTrend(dbClient.QueryExecutor, organizerId)
	if err != nil {
		return nil, err
	}

	var trendEntries []*models.ChartTrend
	for _, trendEntry := range participantsJoiningTrend.Trend {
		trendEntries = append(trendEntries, converters.ConvertTrendEntryRepoToSchema(&trendEntry))
	}
	return &models.ParticipantsJoiningTrend{
		Trend:       trendEntries,
		Granularity: models.Granularity(participantsJoiningTrend.Granularity),
	}, nil
}
