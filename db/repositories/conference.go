package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/maciejas22/conference-manager/api/db"
	filters "github.com/maciejas22/conference-manager/api/db/repositories/shared"
)

type Conference struct {
	Id                   string  `json:"id" db:"id"`
	Title                string  `json:"title" db:"title"`
	Date                 string  `json:"date" db:"date"`
	Location             string  `json:"location" db:"location"`
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

type ConferenceRepository interface {
	GetConference(conferenceId string) (Conference, error)
	GetConferences(page filters.Page, sort *filters.Sort, filters *ConferenceFilter) ([]Conference, filters.PaginationMeta, error)
	CreateConference(conference Conference, organizerId string) (Conference, error)
	UpdateConference(conference Conference) (Conference, error)
}

type conferenceRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewConferenceRepository(ctx context.Context, db *db.DB) ConferenceRepository {
	return &conferenceRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *conferenceRepository) GetConference(conferenceId string) (Conference, error) {
	var conference Conference
	query := "SELECT id, title, date, location, additional_info, participants_limit, registration_deadline FROM " + conference.TableName() + " WHERE id = $1"
	err := r.db.SqlConn.Get(
		&conference,
		query,
		conferenceId,
	)
	if err != nil {
		return Conference{}, err
	}
	return conference, nil
}

func (r *conferenceRepository) GetConferences(p filters.Page, s *filters.Sort, f *ConferenceFilter) ([]Conference, filters.PaginationMeta, error) {
	var conferences []Conference
	var totalItems int
	c := &Conference{}

	query := "SELECT id, title, date, location, additional_info, participants_limit, registration_deadline FROM " + c.TableName()
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

	err := r.db.SqlConn.Get(&totalItems, countQuery, queryArgs[:len(queryArgs)-2]...)
	if err != nil {
		return nil, filters.PaginationMeta{}, err
	}

	err = r.db.SqlConn.Select(&conferences, query, queryArgs...)
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

func (r *conferenceRepository) CreateConference(conference Conference, organizerId string) (Conference, error) {
	query := "INSERT INTO " + conference.TableName() + " (title, date, location, additional_info, participants_limit, registration_deadline) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := r.db.SqlConn.QueryRow(
		query,
		conference.Title,
		conference.Date,
		conference.Location,
		conference.AdditionalInfo,
		conference.ParticipantsLimit,
		conference.RegistrationDeadline,
	).Scan(&conference.Id)
	if err != nil {
		return Conference{}, err
	}

	conferenceOrganizerRepo := NewConferenceOrganizerRepository(r.ctx, r.db)
	_, err = conferenceOrganizerRepo.AddOrganizer(conference.Id, organizerId)
	if err != nil {
		log.Println("Error adding organizer: ", err)
		return Conference{}, errors.New("could not add organizer")
	}
	return conference, nil
}

func (r *conferenceRepository) UpdateConference(conference Conference) (Conference, error) {
	query := "UPDATE " + conference.TableName() + " SET title = $1, date = $2, location = $3, additional_info = $4, participants_limit = $5, registration_deadline = $6 WHERE id = $7"
	_, err := r.db.SqlConn.Exec(
		query,
		conference.Title,
		conference.Date,
		conference.Location,
		conference.AdditionalInfo,
		conference.ParticipantsLimit,
		conference.RegistrationDeadline,
		conference.Id,
	)
	if err != nil {
		return Conference{}, err
	}
	return conference, nil
}
