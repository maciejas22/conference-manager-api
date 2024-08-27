package repositories

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type ConferenceParticipant struct {
	UserId       string `json:"user_id" db:"user_id"`
	ConferenceId string `json:"conference_id" db:"conference_id"`
	JoinedAt     string `json:"joined_at" db:"joined_at"`
}

func (c *ConferenceParticipant) TableName() string {
	return "public.conference_participants"
}

func GetConferenceParticipants(tx *sqlx.Tx, conferenceId string) ([]ConferenceParticipant, error) {
	var participants []ConferenceParticipant
	p := &ConferenceParticipant{}
	query := "SELECT user_id, conference_id FROM " + p.TableName() + " WHERE conference_id = $1"
	err := tx.Select(
		&participants,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func GetConferenceParticipantsCount(tx *sqlx.Tx, conferenceId string) (int, error) {
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

func AddConferenceParticipant(tx *sqlx.Tx, conferenceId string, userId string) (Conference, error) {
	c, err := GetConference(tx, conferenceId)
	if err != nil {
		log.Println("Error getting conference: ", err)
		return Conference{}, errors.New("could not find conference")
	}

	p := &ConferenceParticipant{}
	query := "INSERT INTO " + p.TableName() + " (user_id, conference_id) VALUES ($1, $2)"
	_, err = tx.Exec(query, userId, conferenceId)
	if err != nil {
		log.Println("Error inserting participant: ", err)
		return Conference{}, errors.New("could not insert participant")
	}

	return c, nil
}

func RemoveConferenceParticipant(tx *sqlx.Tx, conferenceId string, userId string) (Conference, error) {
	c, err := GetConference(tx, conferenceId)
	if err != nil {
		log.Println("Error getting conference: ", err)
		return Conference{}, errors.New("could not find conference")
	}

	p := &ConferenceParticipant{}
	query := "DELETE FROM " + p.TableName() + " WHERE user_id = $1 AND conference_id = $2"
	_, err = tx.Exec(query, userId, conferenceId)
	if err != nil {
		log.Println("Error deleting participant: ", err)
		return Conference{}, errors.New("could not delete participant")
	}

	return c, nil
}

func IsConferenceParticipant(tx *sqlx.Tx, conferenceId string, userId string) (bool, error) {
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
