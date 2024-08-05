package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetParticipantsCount(ctx context.Context, db *db.DB, conferenceId string) (int, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	participantsCount, err := repositories.GetConferenceParticipantsCount(tx, conferenceId)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return participantsCount, nil
}

func AddUserToConference(ctx context.Context, db *db.DB, userId string, conferenceID string) (*models.Conference, error) {

	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	conference, err := repositories.AddConferenceParticipant(tx, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
}

func RemoveUserFromConference(ctx context.Context, db *db.DB, userId string, conferenceID string) (*models.Conference, error) {

	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	conference, err := repositories.RemoveConferenceParticipant(tx, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
}

func IsConferenceParticipant(ctx context.Context, db *db.DB, userId, conferenceID string) (*bool, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	isParticipant, err := repositories.IsConferenceOrganizer(tx, userId, conferenceID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &isParticipant, nil
}

func IsConferenceOrganizer(ctx context.Context, db *db.DB, userId string, conferenceID string) (*bool, error) {

	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	isOrganizer, err := repositories.IsConferenceOrganizer(tx, conferenceID, userId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &isOrganizer, nil
}
