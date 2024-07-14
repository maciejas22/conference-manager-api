package repositories

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
)

type ConferenceOrganizer struct {
	UserId       string `json:"user_id" db:"user_id"`
	ConferenceId string `json:"conference_id" db:"conference_id"`
}

func (c *ConferenceOrganizer) TableName() string {
	return "public.conference_organizers"
}

type ConferenceOrganizerRepository interface {
	GetOrganizers(conferenceId string) ([]ConferenceOrganizer, error)
	IsOrganizer(conferenceId string, userId string) (bool, error)
	AddOrganizer(conferenceId string, userId string) (Conference, error)
}

type conferenceOrganizerRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewConferenceOrganizerRepository(ctx context.Context, db *db.DB) ConferenceOrganizerRepository {
	return &conferenceOrganizerRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *conferenceOrganizerRepository) GetOrganizers(conferenceId string) ([]ConferenceOrganizer, error) {
	var organizers []ConferenceOrganizer
	o := &ConferenceOrganizer{}
	query := "SELECT user_id, conference_id FROM " + o.TableName() + " WHERE conference_id = $1"
	err := r.db.SqlConn.Select(
		&organizers,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return organizers, nil
}

func (r *conferenceOrganizerRepository) IsOrganizer(conferenceId string, userId string) (bool, error) {
	var count int
	o := &ConferenceOrganizer{}
	query := "SELECT COUNT(*) FROM " + o.TableName() + " WHERE conference_id = $1 AND user_id = $2"
	err := r.db.SqlConn.Get(
		&count,
		query,
		conferenceId,
		userId,
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *conferenceOrganizerRepository) AddOrganizer(conferenceId string, userId string) (Conference, error) {
	o := &ConferenceOrganizer{
		UserId:       userId,
		ConferenceId: conferenceId,
	}
	query := "INSERT INTO " + o.TableName() + " (user_id, conference_id) VALUES ($1, $2)"
	_, err := r.db.SqlConn.Exec(
		query,
		o.UserId,
		o.ConferenceId,
	)
	if err != nil {
		return Conference{}, err
	}
	return Conference{}, nil
}
