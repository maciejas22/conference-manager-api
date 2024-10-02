package services

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GetParticipantsCount(ctx context.Context, dbClient *db.DB, conferenceId int) (int, error) {
	var participantsCount int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		participantsCount, err = repositories.GetConferenceParticipantsCount(tx, conferenceId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return participantsCount, nil
}

func AddUserToConference(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*int, error) {
	var cId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error

		participantsCount, err := repositories.GetConferenceParticipantsCount(tx, conferenceID)
		if err != nil {
			return err
		}
		conference, err := repositories.GetConference(tx, conferenceID)
		if err != nil {
			return err
		}

		if conference.ParticipantsLimit != nil && participantsCount >= *conference.ParticipantsLimit {
			return gqlerror.Errorf("Conference is full")
		}

		if conference.RegistrationDeadline != nil {
			registrationDeadline, err := time.Parse(time.RFC3339, *conference.RegistrationDeadline)
			if err != nil {
				return err
			}

			if time.Now().After(registrationDeadline) {
				return gqlerror.Errorf("Registration deadline has passed")
			}
		}

		cId, err = repositories.AddConferenceParticipant(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if _, ok := err.(*gqlerror.Error); !ok {
			return nil, gqlerror.Errorf("Failed to add user to conference")
		}
		return nil, err
	}

	return &cId, nil
}

func RemoveUserFromConference(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*int, error) {
	var cId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		cId, err = repositories.RemoveConferenceParticipant(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &cId, nil
}

func IsConferenceParticipant(ctx context.Context, dbClient *db.DB, userId, conferenceID int) (*bool, error) {
	var isParticipant bool
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		isParticipant, err = repositories.IsConferenceParticipant(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &isParticipant, nil
}

func IsConferenceOrganizer(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*bool, error) {
	var isOrganizer bool
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		isOrganizer, err = repositories.IsConferenceOrganizer(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &isOrganizer, nil
}

func GetOrganizerMetrics(ctx context.Context, dbClient *db.DB, organizerId int) (*models.OrganizerMetrics, error) {
	var organizerMetrics repositories.ConferenceOrganizerMetrics
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		organizerMetrics, err = repositories.GetOrganizerLevelMetrics(tx, organizerId)
		if err != nil {
			return err
		}

		return nil
	})
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

func GetParticipantsJoiningTrend(ctx context.Context, dbClient *db.DB, organizerId int) ([]*models.NewParticipantsTrend, error) {
	var participantsJoiningTrend []repositories.TrendEntry
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		participantsJoiningTrend, err = repositories.GetParticipantsTrend(tx, organizerId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var trendEntries []*models.NewParticipantsTrend
	for _, trendEntry := range participantsJoiningTrend {
		date, err := time.Parse(time.RFC3339, trendEntry.Date)
		if err != nil {
			return nil, err
		}

		trendEntries = append(trendEntries, &models.NewParticipantsTrend{
			Date:            date,
			NewParticipants: trendEntry.Value,
		})
	}

	return trendEntries, nil
}
