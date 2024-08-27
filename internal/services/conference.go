package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/internal/utils"
)

func CreateConference(ctx context.Context, db *db.DB, userId string, createConferenceInput models.CreateConferenceInput) (*models.Conference, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	startDate := createConferenceInput.StartDate
	startDateString := utils.TimeToString(&startDate)
	endDate := createConferenceInput.EndDate
	endDateString := utils.TimeToString(&endDate)
	deadlineString := utils.TimeToString(createConferenceInput.RegistrationDeadline)
	conference, err := repositories.CreateConference(tx, repositories.Conference{
		Title:                createConferenceInput.Title,
		StartDate:            *startDateString,
		EndDate:              *endDateString,
		Location:             createConferenceInput.Location,
		Website:              createConferenceInput.Website,
		Acronym:              createConferenceInput.Acronym,
		AdditionalInfo:       createConferenceInput.AdditionalInfo,
		ParticipantsLimit:    createConferenceInput.ParticipantsLimit,
		RegistrationDeadline: deadlineString,
	}, userId)
	if err != nil {
		return nil, err
	}

	for _, a := range createConferenceInput.Agenda {
		startTime := utils.TimeToString(&a.StartTime)
		endTime := utils.TimeToString(&a.EndTime)
		_, err = repositories.CreateAgenda(tx, repositories.AgendaItem{
			ConferenceId: conference.Id,
			StartTime:    *startTime,
			EndTime:      *endTime,
			Event:        a.Event,
			Speaker:      a.Speaker,
		})

		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
}

func ModifyConference(ctx context.Context, db *db.DB, input models.ModifyConferenceInput) (*models.Conference, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	startDate := input.StartDate
	startDateString := utils.TimeToString(startDate)
	endDate := input.EndDate
	endDateString := utils.TimeToString(endDate)
	deadlineString := utils.TimeToString(input.RegistrationDeadline)
	conference, err := repositories.UpdateConference(tx, repositories.Conference{
		Id:                   input.ID,
		Title:                *input.Title,
		StartDate:            *startDateString,
		EndDate:              *endDateString,
		Location:             *input.Location,
		Website:              input.Website,
		Acronym:              input.Acronym,
		AdditionalInfo:       input.AdditionalInfo,
		ParticipantsLimit:    input.ParticipantsLimit,
		RegistrationDeadline: deadlineString,
	})
	if err != nil {
		return nil, err
	}

	for _, a := range input.Agenda {
		startTime := utils.TimeToString(&a.StartTime)
		endTime := utils.TimeToString(&a.EndTime)

		if a.ID != nil && a.Destroy != nil && *a.Destroy {
			err = repositories.DeleteAgenda(tx, *a.ID)
			if err != nil {
				return nil, err
			}
		} else if a.ID == nil {
			_, err = repositories.CreateAgenda(tx, repositories.AgendaItem{
				ConferenceId: conference.Id,
				StartTime:    *startTime,
				EndTime:      *endTime,
				Event:        a.Event,
				Speaker:      a.Speaker,
			})

			if err != nil {
				return nil, err
			}

		} else if a.ID != nil {
			_, err = repositories.UpdateAgenda(tx, repositories.AgendaItem{
				Id:           *a.ID,
				ConferenceId: conference.Id,
				StartTime:    *startTime,
				EndTime:      *endTime,
				Event:        a.Event,
				Speaker:      a.Speaker,
			})

			if err != nil {
				return nil, err
			}

		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
}

func GetAllConferences(ctx context.Context, db *db.DB, userId string, page *models.Page, sort *models.Sort, filters *models.ConferenceFilter) (*models.ConferencePage, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	c, m, err := repositories.GetAllConferences(
		tx,
		userId,
		converters.ConvertPageSchemaToRepo(page),
		converters.ConvertSortSchemaToRepo(sort),
		converters.ConvertConferenceFiltersSchemaToRepo(filters),
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var conferences []*models.Conference
	for _, conference := range c {
		conferences = append(conferences, converters.ConvertConferenceRepoToSchema(&conference))
	}

	return &models.ConferencePage{
		Data: conferences,
		Meta: &models.ConferenceMeta{
			Page: &models.PageInfo{
				TotalItems: m.TotalItems,
				TotalPages: m.TotalPages,
				Number:     m.PageNumber,
				Size:       m.PageSize,
			},
		},
	}, nil
}

func GetConference(ctx context.Context, db *db.DB, id string) (*models.Conference, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	conference, err := repositories.GetConference(tx, id)
	if err != nil {
		return &models.Conference{}, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
}

func GetConferencesMetrics(ctx context.Context, db *db.DB) (*models.ConferencesMetrics, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	metrics, err := repositories.GetMetrics(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.ConferencesMetrics{
		RunningConferences:        metrics.RunningConferences,
		StartingInLessThan24Hours: metrics.StartingInLessThan24Hours,
		TotalConducted:            metrics.TotalConducted,
		ParticipantsToday:         metrics.ParticipantsToday,
	}, nil
}
