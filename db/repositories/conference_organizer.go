package repositories

import (
	"github.com/jmoiron/sqlx"
)

type ConferenceOrganizer struct {
	UserId       string `json:"user_id" db:"user_id"`
	ConferenceId string `json:"conference_id" db:"conference_id"`
}

func (c *ConferenceOrganizer) TableName() string {
	return "public.conference_organizers"
}

func GetConferenceOrganizers(tx *sqlx.Tx, conferenceId string) ([]ConferenceOrganizer, error) {
	var organizers []ConferenceOrganizer
	o := &ConferenceOrganizer{}
	query := "SELECT user_id, conference_id FROM " + o.TableName() + " WHERE conference_id = $1"
	err := tx.Select(
		&organizers,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return organizers, nil
}

func IsConferenceOrganizer(tx *sqlx.Tx, conferenceId string, userId string) (bool, error) {
	var count int
	o := &ConferenceOrganizer{}
	query := "SELECT COUNT(*) FROM " + o.TableName() + " WHERE conference_id = $1 AND user_id = $2"
	err := tx.Get(
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

func AddConferenceOrganizer(tx *sqlx.Tx, conferenceId string, userId string) (Conference, error) {
	o := &ConferenceOrganizer{
		UserId:       userId,
		ConferenceId: conferenceId,
	}
	query := "INSERT INTO " + o.TableName() + " (user_id, conference_id) VALUES ($1, $2)"
	_, err := tx.Exec(
		query,
		o.UserId,
		o.ConferenceId,
	)
	if err != nil {
		return Conference{}, err
	}
	return Conference{}, nil
}
