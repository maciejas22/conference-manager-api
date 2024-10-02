package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type ConferenceParticipant struct {
	UserId       int    `json:"user_id" db:"user_id"`
	ConferenceId int    `json:"conference_id" db:"conference_id"`
	JoinedAt     string `json:"joined_at" db:"joined_at"`
}

func (c *ConferenceParticipant) TableName() string {
	return "public.conference_participants"
}

func GetConferenceParticipantsCount(tx *sqlx.Tx, conferenceId int) (int, error) {
	var count int
	p := &ConferenceParticipant{}
	query := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE conference_id = $1"
	err := tx.Get(
		&count,
		query,
		conferenceId,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddConferenceParticipant(tx *sqlx.Tx, conferenceId int, userId int) (int, error) {
	p := &ConferenceParticipant{}
	query := "INSERT INTO " + p.TableName() + " (user_id, conference_id) VALUES ($1, $2)"
	_, err := tx.Exec(query, userId, conferenceId)
	if err != nil {
		return 0, errors.New("could not insert participant")
	}

	return conferenceId, nil
}

func RemoveConferenceParticipant(tx *sqlx.Tx, conferenceId int, userId int) (int, error) {
	p := &ConferenceParticipant{}
	query := "DELETE FROM " + p.TableName() + " WHERE user_id = $1 AND conference_id = $2"
	_, err := tx.Exec(query, userId, conferenceId)
	if err != nil {
		return 0, errors.New("could not delete participant")
	}

	return conferenceId, nil
}

func IsConferenceParticipant(tx *sqlx.Tx, conferenceId int, userId int) (bool, error) {
	var count int
	p := &ConferenceParticipant{}
	query := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE conference_id = $1 AND user_id = $2"
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
