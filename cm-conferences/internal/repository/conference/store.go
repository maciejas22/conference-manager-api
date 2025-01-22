package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"
)

type ConferenceRepo struct {
	Db *sqlx.DB
}

func NewConferenceRepo(db *sqlx.DB) ConferenceRepoInterface {
	return &ConferenceRepo{Db: db}
}

type ConferenceFilter struct {
	Title          *string
	AssociatedOnly *bool
	RunningOnly    *bool
}

func (r *ConferenceRepo) GetConferenceById(conferenceId int) (Conference, error) {
	var conference Conference
	query := "SELECT id, title, start_date, end_date, location, additional_info, acronym, website, participants_limit, registration_deadline, ticket_price FROM " + common.GetTableName(common.ConferenceTable) + " WHERE id = $1"
	err := r.Db.Get(
		&conference,
		query,
		conferenceId,
	)
	if err != nil {
		return Conference{}, errors.New("Conference not found")
	}
	return conference, nil
}

func (r *ConferenceRepo) GetConferencesByIds(conferenceIds []int) ([]Conference, error) {
	var conferences []Conference
	query := "SELECT id, title, start_date, end_date, location, additional_info, acronym, website, participants_limit, registration_deadline, ticket_price FROM " + common.GetTableName(common.ConferenceTable) + " WHERE id = ANY($1)"
	err := r.Db.Select(&conferences, query, conferenceIds)
	if err != nil {
		return nil, errors.New("Conferences not found")
	}
	return conferences, nil
}

func (r *ConferenceRepo) GetConferencesPage(userId int, p common.Page, s *common.Sort, f *ConferenceFilter) ([]int, common.PaginationMeta, error) {
	var conferenceIds []int
	var totalItems int

	query := "SELECT id FROM " + common.GetTableName(common.ConferenceTable)
	countQuery := "SELECT COUNT(*) FROM " + common.GetTableName(common.ConferenceTable)

	whereClause := " WHERE 1=1"
	queryArgs := []interface{}{}

	if f != nil && f.Title != nil {
		whereClause += fmt.Sprintf(" AND title ILIKE $%d", len(queryArgs)+1)
		queryArgs = append(queryArgs, "%"+*f.Title+"%")
	}

	if f != nil && f.AssociatedOnly != nil && *f.AssociatedOnly {
		whereClause += fmt.Sprintf(" AND id IN (SELECT conference_id FROM %s WHERE user_id = $%d UNION SELECT conference_id FROM %s WHERE user_id = $%d)",
			common.GetTableName(common.ConferenceOrganizerTable), len(queryArgs)+1, common.GetTableName(common.ConferenceParticipantTable), len(queryArgs)+2)

		queryArgs = append(queryArgs, userId, userId)
	}

	if f != nil && f.RunningOnly != nil && *f.RunningOnly {
		now := time.Now()
		whereClause += fmt.Sprintf(" AND start_date <= $%d AND end_date >= $%d", len(queryArgs)+1, len(queryArgs)+2)
		queryArgs = append(queryArgs, now, now)
	}

	query += whereClause
	countQuery += whereClause

	if s != nil {
		query += fmt.Sprintf(" ORDER BY %s %s", s.Column, s.Order)
	}

	offset := (p.PageNumber - 1) * p.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(queryArgs)+1, len(queryArgs)+2)
	queryArgs = append(queryArgs, p.PageSize, offset)

	err := r.Db.Get(&totalItems, countQuery, queryArgs[:len(queryArgs)-2]...)
	if err != nil {
		return nil, common.PaginationMeta{}, errors.New("Could not count conferences")
	}

	err = r.Db.Select(&conferenceIds, query, queryArgs...)
	if err != nil {
		return nil, common.PaginationMeta{}, errors.New("Could not get conferences")
	}

	totalPages := (totalItems + p.PageSize - 1) / p.PageSize

	paginationMeta := common.PaginationMeta{
		TotalItems: totalItems,
		TotalPages: totalPages,
		PageSize:   p.PageSize,
		PageNumber: p.PageNumber,
	}
	return conferenceIds, paginationMeta, nil
}

func (r *ConferenceRepo) CreateConference(conference Conference, organizerId int) (int, error) {
	query := "INSERT INTO " + common.GetTableName(common.ConferenceTable) + " (title, start_date, end_date, location, website, acronym, additional_info, participants_limit, registration_deadline, ticket_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	var conferenceId int
	err := r.Db.QueryRowx(query,
		conference.Title,
		conference.StartDate,
		conference.EndDate,
		conference.Location,
		conference.Website,
		conference.Acronym,
		conference.AdditionalInfo,
		conference.ParticipantsLimit,
		conference.RegistrationDeadline,
		conference.TicketPrice,
	).Scan(&conferenceId)
	if err != nil {
		return 0, errors.New("Could not create conference")
	}

	return conferenceId, nil
}

func (r *ConferenceRepo) UpdateConference(conference Conference) (int, error) {
	query := "UPDATE " + common.GetTableName(common.ConferenceTable) + " SET title = $1, start_date = $2, end_date = $3, location = $4, website = $5, acronym = $6, additional_info = $7, participants_limit = $8, registration_deadline = $9, ticket_price = $10 WHERE id = $11"
	_, err := r.Db.Exec(
		query,
		conference.Title,
		conference.StartDate,
		conference.EndDate,
		conference.Location,
		conference.Website,
		conference.Acronym,
		conference.AdditionalInfo,
		conference.ParticipantsLimit,
		conference.RegistrationDeadline,
		conference.TicketPrice,
		conference.Id,
	)
	if err != nil {
		return 0, errors.New("Could not update conference")
	}
	return conference.Id, nil
}

type ConferencesMetrics struct {
	RunningConferences        int `json:"runningConferences"`
	StartingInLessThan24Hours int `json:"startingInLessThan24Hours"`
	TotalConducted            int `json:"totalConducted"`
	ParticipantsToday         int `json:"participantsToday"`
}

func (r *ConferenceRepo) GetMetrics() (ConferencesMetrics, error) {
	var metrics ConferencesMetrics

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	err := r.Db.Get(&metrics.RunningConferences, "SELECT COUNT(*) FROM "+common.GetTableName(common.ConferenceTable)+" WHERE start_date <= $1 AND end_date >= $2", now, now)
	if err != nil {
		return ConferencesMetrics{}, errors.New("Could not count running conferences")
	}

	err = r.Db.Get(&metrics.StartingInLessThan24Hours, "SELECT COUNT(*) FROM "+common.GetTableName(common.ConferenceTable)+" WHERE start_date >= $1 AND start_date <= $2", now, tomorrow)
	if err != nil {
		return ConferencesMetrics{}, errors.New("Could not count conferences starting in less than 24 hours")
	}

	err = r.Db.Get(&metrics.TotalConducted, "SELECT COUNT(*) FROM "+common.GetTableName(common.ConferenceTable))
	if err != nil {
		return ConferencesMetrics{}, errors.New("Could not count total conducted conferences")
	}

	err = r.Db.Get(&metrics.ParticipantsToday, "SELECT COUNT(DISTINCT user_id) FROM "+common.GetTableName(common.ConferenceParticipantTable)+" cp JOIN "+common.GetTableName(common.ConferenceTable)+" c ON cp.conference_id = c.id WHERE c.start_date <= $1 AND c.end_date >= $2", now, now)
	if err != nil {
		return ConferencesMetrics{}, errors.New("Could not count participants today")
	}

	return metrics, nil
}
