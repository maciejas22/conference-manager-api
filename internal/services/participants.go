package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	filters "github.com/maciejas22/conference-manager/api/internal/db/repositories/shared"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GetParticipantsCount(ctx context.Context, dbClient *db.DB, conferenceId int) (int, error) {
	var participantsCount int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		participantsCount, err = repositories.GetConferenceParticipantsCount(tx, conferenceId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return participantsCount, nil
}

func AddUserToConference(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*int, error) {
	var cId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error

		cId, err = repositories.AddConferenceParticipant(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if _, ok := err.(*gqlerror.Error); !ok {
			return nil, gqlerror.Errorf("Failed to add user to conference")
		}
		return nil, err
	}

	return &cId, nil
}

func RemoveUserFromConference(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*int, error) {
	var cId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		cId, err = repositories.RemoveConferenceParticipant(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &cId, nil
}

func IsConferenceParticipant(ctx context.Context, dbClient *db.DB, userId, conferenceID int) (*bool, error) {
	var isParticipant bool
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		isParticipant, err = repositories.IsConferenceParticipant(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &isParticipant, nil
}

func IsConferenceOrganizer(ctx context.Context, dbClient *db.DB, userId int, conferenceID int) (*bool, error) {
	var isOrganizer bool
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		isOrganizer, err = repositories.IsConferenceOrganizer(tx, conferenceID, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &isOrganizer, nil
}

func GetOrganizerMetrics(ctx context.Context, dbClient *db.DB, organizerId int) (*models.OrganizerMetrics, error) {
	var organizerMetrics repositories.ConferenceOrganizerMetrics
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		organizerMetrics, err = repositories.GetOrganizerLevelMetrics(tx, organizerId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &models.OrganizerMetrics{
		RunningConferences:        organizerMetrics.RunningConferencesCount,
		ParticipantsCount:         organizerMetrics.ParticipantsCount,
		AverageParticipantsCount:  organizerMetrics.AverageParticipantsCount,
		TotalOrganizedConferences: organizerMetrics.TotalOrganizedConferences,
	}, nil
}

func GetParticipantsJoiningTrend(ctx context.Context, dbClient *db.DB, organizerId int) ([]*models.NewParticipantsTrend, error) {
	var participantsJoiningTrend []repositories.TrendEntry
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		participantsJoiningTrend, err = repositories.GetParticipantsTrend(tx, organizerId)
		if err != nil {
			return errors.New("Failed to fetch participants joining trend")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var trendEntries []*models.NewParticipantsTrend
	for _, trendEntry := range participantsJoiningTrend {
		date, err := time.Parse(time.RFC3339, trendEntry.Date)
		if err != nil {
			return nil, errors.New("Failed to parse date")
		}

		trendEntries = append(trendEntries, &models.NewParticipantsTrend{
			Date:            date,
			NewParticipants: trendEntry.Value,
		})
	}

	return trendEntries, nil
}

func GetParticipantsTickets(ctx context.Context, dbClient *db.DB, participantId int, page filters.Page) ([]*models.Ticket, *models.PageInfo, error) {
	var ticketData []*models.Ticket
	var ticketsMeta *models.PageInfo
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		participantTickets, participantsmeta, err := repositories.GetParticipantTickets(tx, participantId, page)
		if err != nil {
			return err
		}

		ticketsMeta = &models.PageInfo{
			TotalItems: participantsmeta.TotalItems,
			TotalPages: participantsmeta.TotalPages,
			Size:       participantsmeta.PageSize,
			Number:     participantsmeta.PageNumber,
		}

		conferenceIds := make([]int, len(participantTickets))
		for i, ticket := range participantTickets {
			conferenceIds[i] = ticket.ConferenceId
		}

		conferences, err := repositories.GetConferencesByIds(tx, conferenceIds)
		if err != nil {
			return err
		}

		conferenceMap := make(map[int]*models.Conference)
		for _, conference := range conferences {
			startDate, err := time.Parse(time.RFC3339, conference.StartDate)
			if err != nil {
				return err
			}

			endDate, err := time.Parse(time.RFC3339, conference.EndDate)
			if err != nil {
				return err
			}

			var registrationDeadline *time.Time
			if conference.RegistrationDeadline != nil {
				parsedDeadline, err := time.Parse(time.RFC3339, *conference.RegistrationDeadline)
				if err != nil {
					return err
				}
				registrationDeadline = &parsedDeadline
			}

			conferenceMap[conference.Id] = &models.Conference{
				ID:                   conference.Id,
				Title:                conference.Title,
				StartDate:            startDate,
				EndDate:              endDate,
				Location:             conference.Location,
				Website:              conference.Website,
				Acronym:              conference.Acronym,
				AdditionalInfo:       conference.AdditionalInfo,
				ParticipantsLimit:    conference.ParticipantsLimit,
				RegistrationDeadline: registrationDeadline,
				TicketPrice:          conference.TicketPrice,
			}
		}

		for _, ticket := range participantTickets {
			conference := conferenceMap[ticket.ConferenceId]
			if conference == nil {
				return fmt.Errorf("Conference with ID %d not found", ticket.ConferenceId)
			}
			ticketData = append(ticketData, &models.Ticket{
				ID:         ticket.TicketId,
				Conference: conference,
			})
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return ticketData, ticketsMeta, nil
}
