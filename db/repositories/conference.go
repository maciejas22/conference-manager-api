package repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	filters "github.com/maciejas22/conference-manager/api/db/repositories/shared"
)

type Conference struct {
	Id                   string  `json:"id" db:"id"`
	Title                string  `json:"title" db:"title"`
	StartDate            string  `json:"start_date" db:"start_date"`
	EndDate              string  `json:"end_date" db:"end_date"`
	Location             string  `json:"location" db:"location"`
	Website              *string `json:"website,omitempty" db:"website"`
	Acronym              *string `json:"acronym,omitempty" db:"acronym"`
	AdditionalInfo       *string `json:"additional_info,omitempty" db:"additional_info"`
	ParticipantsLimit    *int    `json:"participants_limit,omitempty" db:"participants_limit"`
	RegistrationDeadline *string `json:"registration_deadline,omitempty" db:"registration_deadline"`
}

func (c *Conference) TableName() string {
	return "public.conferences"
}

type ConferenceFilter struct {
	Title          *string `json:"title"`
	AssociatedOnly *bool   `json:"associatedOnly,omitempty"`
}

func GetConference(tx *sqlx.Tx, conferenceId string) (Conference, error) {
	var conference Conference
	query := "SELECT id, title, start_date, end_date, location, additional_info, acronym, website, participants_limit, registration_deadline FROM " + conference.TableName() + " WHERE id = $1"
	err := tx.Get(
		&conference,
		query,
		conferenceId,
	)
	if err != nil {
		return Conference{}, err
	}
	return conference, nil
}

func GetAllConferences(tx *sqlx.Tx, p filters.Page, s *filters.Sort, f *ConferenceFilter) ([]Conference, filters.PaginationMeta, error) {
	var conferences []Conference
	var totalItems int
	c := &Conference{}

	query := "SELECT id, title, start_date, end_date, location, website, acronym, additional_info, participants_limit, registration_deadline FROM " + c.TableName()
	countQuery := "SELECT COUNT(*) FROM " + c.TableName()

	whereClause := " WHERE 1=1"
	queryArgs := []interface{}{}
	argCounter := 1

	if f.Title != nil {
		whereClause += fmt.Sprintf(" AND title LIKE $%d", argCounter)
		queryArgs = append(queryArgs, "%"+*f.Title+"%")
		argCounter++
	}

	query += whereClause
	countQuery += whereClause

	if s != nil {
		query += fmt.Sprintf(" ORDER BY %s %s", s.Column, s.Order)
	}

	offset := (p.PageNumber - 1) * p.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	queryArgs = append(queryArgs, p.PageSize, offset)

	err := tx.Get(&totalItems, countQuery, queryArgs[:len(queryArgs)-2]...)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	err = tx.Select(&conferences, query, queryArgs...)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	totalPages := (totalItems + p.PageSize - 1) / p.PageSize

	paginationMeta := filters.PaginationMeta{
		TotalItems: totalItems,
		TotalPages: totalPages,
		PageSize:   p.PageSize,
		PageNumber: p.PageNumber,
	}
	return conferences, paginationMeta, nil
}

func CreateConference(tx *sqlx.Tx, conference Conference, organizerId string) (Conference, error) {
	query := "INSERT INTO " + conference.TableName() + " (title, start_date, end_date, location, website, acronym, additional_info, participants_limit, registration_deadline) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := tx.Get(&conference.Id,
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
	)
	if err != nil {
		return Conference{}, err
	}

	_, err = AddConferenceOrganizer(tx, conference.Id, organizerId)
	if err != nil {
		log.Println("Error adding organizer: ", err)
		return Conference{}, errors.New("could not add organizer")
	}
	return conference, nil
}

func UpdateConference(tx *sqlx.Tx, conference Conference) (Conference, error) {
	query := "UPDATE " + conference.TableName() + " SET title = $1, start_date = $2, end_date = $3, location = $4, website = $5, acronym = $6, additional_info = $7, participants_limit = $8, registration_deadline = $9 WHERE id = $10"
	_, err := tx.Exec(
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
		conference.Id,
	)
	if err != nil {
		log.Println("Error updating conference: ", err)
		return Conference{}, err
	}
	return conference, nil
}

type ConferencesMetrics struct {
	RunningConferences        int `json:"runningConferences"`
	StartingInLessThan24Hours int `json:"startingInLessThan24Hours"`
	TotalConducted            int `json:"totalConducted"`
	ParticipantsToday         int `json:"participantsToday"`
}

func GetMetrics(tx *sqlx.Tx) (ConferencesMetrics, error) {
	var conference Conference
	var metrics ConferencesMetrics

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	err := tx.Get(&metrics.RunningConferences, "SELECT COUNT(*) FROM "+conference.TableName()+" WHERE start_date <= $1 AND end_date >= $1", now)
	if err != nil {
		return ConferencesMetrics{}, err
	}

	err = tx.Get(&metrics.StartingInLessThan24Hours, "SELECT COUNT(*) FROM "+conference.TableName()+" WHERE start_date >= $1 AND start_date <= $2", now, tomorrow)
	if err != nil {
		return ConferencesMetrics{}, err
	}

	err = tx.Get(&metrics.TotalConducted, "SELECT COUNT(*) FROM "+conference.TableName())
	if err != nil {
		return ConferencesMetrics{}, err
	}

	var conferenceParticipant ConferenceParticipant
	err = tx.Get(&metrics.ParticipantsToday, "SELECT COUNT(DISTINCT user_id) FROM "+conferenceParticipant.TableName()+" cp JOIN "+conference.TableName()+" c ON cp.conference_id = c.id WHERE c.start_date <= $1 AND c.end_date >= $1", now)
	if err != nil {
		return ConferencesMetrics{}, err
	}

	return metrics, nil
}
