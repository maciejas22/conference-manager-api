package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
)

type ConferenceParticipant struct {
	UserId       int    `json:"user_id" db:"user_id"`
	ConferenceId int    `json:"conference_id" db:"conference_id"`
	JoinedAt     string `json:"joined_at" db:"joined_at"`
}

func (c *ConferenceParticipant) TableName() string {
	return "conference_participants"
}

func GetConferenceParticipants(qe *db.QueryExecutor, conferenceId int) ([]ConferenceParticipant, error) {
	var participants []ConferenceParticipant
	p := &ConferenceParticipant{}
	query := "SELECT user_id, conference_id FROM " + p.TableName() + " WHERE conference_id = ?"
	err := sqlx.Select(
		qe,
		&participants,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func GetConferenceParticipantsCount(qe *db.QueryExecutor, conferenceId int) (int, error) {
	var count int
	p := &ConferenceParticipant{}
	query := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE conference_id = ?"
	err := sqlx.Get(
		qe,
		&count,
		query,
		conferenceId,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddConferenceParticipant(qe *db.QueryExecutor, conferenceId int, userId int) (Conference, error) {
	c, err := GetConference(qe, conferenceId)
	if err != nil {
		return Conference{}, errors.New("could not find conference")
	}

	p := &ConferenceParticipant{}
	query := "INSERT INTO " + p.TableName() + " (user_id, conference_id) VALUES (?, ?)"
	_, err = qe.Exec(query, userId, conferenceId)
	if err != nil {
		return Conference{}, errors.New("could not insert participant")
	}

	return c, nil
}

func RemoveConferenceParticipant(qe *db.QueryExecutor, conferenceId int, userId int) (Conference, error) {
	c, err := GetConference(qe, conferenceId)
	if err != nil {
		return Conference{}, errors.New("could not find conference")
	}

	p := &ConferenceParticipant{}
	query := "DELETE FROM " + p.TableName() + " WHERE user_id = ? AND conference_id = ?"
	_, err = qe.Exec(query, userId, conferenceId)
	if err != nil {
		return Conference{}, errors.New("could not delete participant")
	}

	return c, nil
}

func IsConferenceParticipant(qe *db.QueryExecutor, conferenceId int, userId int) (bool, error) {
	var count int
	p := &ConferenceParticipant{}
	query := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE conference_id = ? AND user_id = ?"
	err := sqlx.Get(
		qe,
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
