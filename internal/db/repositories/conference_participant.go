package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
	filters "github.com/maciejas22/conference-manager/api/internal/db/repositories/shared"
)

type ConferenceParticipant struct {
	UserId       int    `json:"user_id" db:"user_id"`
	ConferenceId int    `json:"conference_id" db:"conference_id"`
	TicketId     string `json:"ticket_id" db:"ticket_id"`
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

func GetParticipantTickets(tx *sqlx.Tx, userId int, page filters.Page) ([]ConferenceParticipant, filters.PaginationMeta, error) {
	var tickets []ConferenceParticipant
	var totalItems int
	p := &ConferenceParticipant{}
	query := "SELECT conference_id, ticket_id, joined_at FROM " + p.TableName() + " WHERE user_id = $1 ORDER BY joined_at DESC LIMIT $2 OFFSET $3"
	countQuery := "SELECT COUNT(*) FROM " + p.TableName() + " WHERE user_id = $1"

	err := tx.Select(
		&tickets,
		query,
		userId,
		page.PageSize,
		(page.PageNumber-1)*page.PageSize,
	)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	err = tx.Get(
		&totalItems,
		countQuery,
		userId,
	)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	return tickets, filters.PaginationMeta{
		TotalItems: totalItems,
		PageSize:   page.PageSize,
		PageNumber: page.PageNumber,
	}, nil
}
