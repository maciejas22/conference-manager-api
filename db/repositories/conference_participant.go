package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/maciejas22/conference-manager/api/db"
)

type ConferenceParticipant struct {
	UserId       string `json:"user_id" db:"user_id"`
	ConferenceId string `json:"conference_id" db:"conference_id"`
}

func (c *ConferenceParticipant) TableName() string {
	return "public.conference_participants"
}

type ConferenceParticipantRepository interface {
	GetParticipants(conferenceId string) ([]ConferenceParticipant, error)
	GetParticipantsCount(conferenceId string) (int, error)
	AddParticipant(conferenceId string, userId string) (Conference, error)
	RemoveParticipant(conferenceId string, userId string) (Conference, error)
	isParticipant(conferenceId string, userId string) (bool, error)
}

type conferenceParticipantRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewConferenceParticipantRepository(ctx context.Context, db *db.DB) ConferenceParticipantRepository {
	return &conferenceParticipantRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *conferenceParticipantRepository) GetParticipants(conferenceId string) ([]ConferenceParticipant, error) {
	var participants []ConferenceParticipant
	p := &ConferenceParticipant{}
	query := "SELECT user_id, conference_id FROM " + p.TableName() + " WHERE conference_id = $1"
	err := r.db.SqlConn.Select(
		&participants,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *conferenceParticipantRepository) GetParticipantsCount(conferenceId string) (int, error) {
	var count int
	p := &ConferenceParticipant{}
	query := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE conference_id = $1"
	err := r.db.SqlConn.Get(
		&count,
		query,
		conferenceId,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *conferenceParticipantRepository) AddParticipant(conferenceId string, userId string) (Conference, error) {
	conferenceRepo := NewConferenceRepository(r.ctx, r.db)
	c, err := conferenceRepo.GetConference(conferenceId)
	if err != nil {
		log.Println("Error getting conference: ", err)
		return Conference{}, errors.New("could not find conference")
	}

	p := &ConferenceParticipant{}
	query := "INSERT INTO " + p.TableName() + " (user_id, conference_id) VALUES ($1, $2)"
	_, err = r.db.SqlConn.Exec(query, userId, conferenceId)
	if err != nil {
		log.Println("Error inserting participant: ", err)
		return Conference{}, errors.New("could not insert participant")
	}

	return c, nil
}

func (r *conferenceParticipantRepository) RemoveParticipant(conferenceId string, userId string) (Conference, error) {
	conferenceRepo := NewConferenceRepository(r.ctx, r.db)
	c, err := conferenceRepo.GetConference(conferenceId)
	if err != nil {
		log.Println("Error getting conference: ", err)
		return Conference{}, errors.New("could not find conference")
	}

	p := &ConferenceParticipant{}
	query := "DELETE FROM " + p.TableName() + " WHERE user_id = $1 AND conference_id = $2"
	_, err = r.db.SqlConn.Exec(query, userId, conferenceId)
	if err != nil {
		log.Println("Error deleting participant: ", err)
		return Conference{}, errors.New("could not delete participant")
	}

	return c, nil
}

func (r *conferenceParticipantRepository) isParticipant(conferenceId string, userId string) (bool, error) {
	var count int
	p := &ConferenceParticipant{}
	query := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE conference_id = $1 AND user_id = $2"
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
