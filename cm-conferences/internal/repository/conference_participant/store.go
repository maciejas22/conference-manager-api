package repository

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"
)

type ConferenceParticipantRepo struct {
	Db *sqlx.DB
}

func NewParticipantRepo(db *sqlx.DB) ConferenceParticipantRepoInterface {
	return &ConferenceParticipantRepo{Db: db}
}

func (r *ConferenceParticipantRepo) GetConferenceParticipantsCount(conferenceIds []int) ([]ParticipantsCount, error) {
	results := []ParticipantsCount{}
	query := "SELECT conference_id, COUNT(*) as count FROM " + common.GetTableName(common.ConferenceParticipantTable) + " WHERE conference_id = ANY($1) GROUP BY conference_id"

	err := r.Db.Select(&results, query, conferenceIds)
	if err != nil {
		return nil, errors.New("Could not get participants count")
	}

	return results, nil
}

func (r *ConferenceParticipantRepo) AddConferenceParticipant(conferenceId int, userId int) (string, error) {
	query := "INSERT INTO " + common.GetTableName(common.ConferenceParticipantTable) + " (user_id, conference_id) VALUES ($1, $2) RETURNING ticket_id"
	var ticketId string
	err := r.Db.QueryRow(query, userId, conferenceId).Scan(&ticketId)
	if err != nil {
		return "", errors.New("Could not insert participant")
	}

	return ticketId, nil
}

func (r *ConferenceParticipantRepo) RemoveConferenceParticipant(conferenceId int, userId int) (int, error) {
	query := "DELETE FROM " + common.GetTableName(common.ConferenceParticipantTable) + " WHERE user_id = $1 AND conference_id = $2"
	_, err := r.Db.Exec(query, userId, conferenceId)
	if err != nil {
		return 0, errors.New("Could not delete participant")
	}

	return conferenceId, nil
}

func (r *ConferenceParticipantRepo) IsConferenceParticipant(conferenceId int, userId int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + common.GetTableName(common.ConferenceParticipantTable) + " WHERE conference_id = $1 AND user_id = $2"
	err := r.Db.Get(
		&count,
		query,
		conferenceId,
		userId,
	)
	if err != nil {
		return false, errors.New("Could not check if user is a participant")
	}
	return count > 0, nil
}

func (r *ConferenceParticipantRepo) GetParticipantTickets(userId int, page common.Page) ([]ConferenceParticipant, common.PaginationMeta, error) {
	var tickets []ConferenceParticipant
	var totalItems int
	query := "SELECT conference_id, ticket_id, joined_at FROM " + common.GetTableName(common.ConferenceParticipantTable) + " WHERE user_id = $1 ORDER BY joined_at DESC LIMIT $2 OFFSET $3"
	countQuery := "SELECT COUNT(*) FROM " + common.GetTableName(common.ConferenceParticipantTable) + " WHERE user_id = $1"

	err := r.Db.Select(
		&tickets,
		query,
		userId,
		page.PageSize,
		(page.PageNumber-1)*page.PageSize,
	)
	if err != nil {
		return nil, common.PaginationMeta{}, errors.New("Could not get tickets")
	}

	err = r.Db.Get(
		&totalItems,
		countQuery,
		userId,
	)
	if err != nil {
		return nil, common.PaginationMeta{}, errors.New("Could not count tickets")
	}

	return tickets, common.PaginationMeta{
		TotalItems: totalItems,
		PageSize:   page.PageSize,
		PageNumber: page.PageNumber,
	}, nil
}

func (r *ConferenceParticipantRepo) IsTicketValid(ticketId string, conferenceId int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + common.GetTableName(common.ConferenceParticipantTable) + " WHERE ticket_id = $1 AND conference_id = $2"
	err := r.Db.Get(
		&count,
		query,
		ticketId,
		conferenceId,
	)
	if err != nil {
		log.Println(err)
		return false, errors.New("Could not check if ticket is valid")
	}
	return count > 0, nil
}
