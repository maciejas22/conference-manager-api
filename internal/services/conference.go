package services

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	filters "github.com/maciejas22/conference-manager/api/internal/db/repositories/shared"
	"github.com/maciejas22/conference-manager/api/internal/files"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/internal/utils"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func CreateConference(ctx context.Context, dbClient *db.DB, s3 *s3.S3Client, userId int, createConferenceInput models.CreateConferenceInput) (*int, error) {
	var conferenceId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		startDate := createConferenceInput.StartDate
		startDateString := utils.TimeToString(&startDate)
		endDate := createConferenceInput.EndDate
		endDateString := utils.TimeToString(&endDate)
		deadlineString := utils.TimeToString(createConferenceInput.RegistrationDeadline)
		conferenceId, err := repositories.CreateConference(tx, repositories.Conference{
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
			return err
		}

		for _, a := range createConferenceInput.Agenda {
			startTime := utils.TimeToString(&a.CreateItem.StartTime)
			endTime := utils.TimeToString(&a.CreateItem.EndTime)
			err = repositories.CreateAgenda(tx, repositories.AgendaItem{
				ConferenceId: conferenceId,
				StartTime:    *startTime,
				EndTime:      *endTime,
				Event:        a.CreateItem.Event,
				Speaker:      a.CreateItem.Speaker,
			})

			if err != nil {
				return err
			}
		}

		for _, f := range createConferenceInput.Files {
			fileParts := strings.Split(f.UploadFile.Base64Content, ".")
			if len(fileParts) < 1 {
				f.UploadFile.Base64Content = fileParts[1]
			}

			fileData, err := base64.StdEncoding.DecodeString(f.UploadFile.Base64Content)
			if err != nil {
				return err
			}
			err = files.UploadConferenceFile(ctx, s3, conferenceId, f.UploadFile.FileName, strings.NewReader(string(fileData)))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return &conferenceId, err
}

func ModifyConference(ctx context.Context, dbClient *db.DB, s3 *s3.S3Client, input models.ModifyConferenceInput) (*int, error) {
	var conferenceId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		startDate := input.StartDate
		startDateString := utils.TimeToString(startDate)
		endDate := input.EndDate
		endDateString := utils.TimeToString(endDate)
		deadlineString := utils.TimeToString(input.RegistrationDeadline)
		conferenceId, err = repositories.UpdateConference(tx, repositories.Conference{
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
			return err
		}

		for _, a := range input.Agenda {
			if a.CreateItem != nil {
				startTime := utils.TimeToString(&a.CreateItem.StartTime)
				endTime := utils.TimeToString(&a.CreateItem.EndTime)
				err = repositories.CreateAgenda(tx, repositories.AgendaItem{
					ConferenceId: input.ID,
					StartTime:    *startTime,
					EndTime:      *endTime,
					Event:        a.CreateItem.Event,
					Speaker:      a.CreateItem.Speaker,
				})
				if err != nil {
					return err
				}
			} else if a.DeleteItem != nil {
				err = repositories.DeleteAgenda(tx, *a.DeleteItem)
				if err != nil {
					return err
				}
			}
		}

		for _, f := range input.Files {
			if f.DeleteFile != nil {
				err = files.DeleteConferenceFile(ctx, s3, f.DeleteFile.Key)
				if err != nil {
					return err
				}
			} else if f.UploadFile != nil {
				fileParts := strings.Split(f.UploadFile.Base64Content, ".")
				if len(fileParts) < 1 {
					f.UploadFile.Base64Content = fileParts[1]
				}

				fileData, err := base64.StdEncoding.DecodeString(f.UploadFile.Base64Content)
				if err != nil {
					return err
				}

				err = files.UploadConferenceFile(ctx, s3, input.ID, f.UploadFile.FileName, strings.NewReader(string(fileData)))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return &conferenceId, err
}

func GetAllConferences(ctx context.Context, dbClient *db.DB, userId int, p *models.Page, s *models.Sort, f *models.ConferencesFilters) (*models.ConferencesPage, error) {
	var c []repositories.Conference
	var m filters.PaginationMeta

	var page filters.Page
	if p == nil {
		page = filters.Page{
			PageNumber: 1,
			PageSize:   10,
		}
	} else {
		page = filters.Page{
			PageNumber: p.Number,
			PageSize:   p.Size,
		}
	}

	var sort *filters.Sort
	if s != nil {
		var order filters.Order
		if s.Order == models.OrderAsc {
			order = filters.ASC
		} else {
			order = filters.DESC
		}

		sort = &filters.Sort{
			Column: s.Column,
			Order:  order,
		}
	}

	var filters *repositories.ConferenceFilter
	if f != nil {
		filters = &repositories.ConferenceFilter{
			Title:          f.Title,
			AssociatedOnly: f.AssociatedOnly,
		}
	}

	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		c, m, err = repositories.GetAllConferences(
			tx,
			userId,
			page,
			sort,
			filters,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return &models.ConferencesPage{}, err
	}

	var conferences []*models.Conference
	for _, conference := range c {
		startDate, err := time.Parse(time.RFC3339, conference.StartDate)
		if err != nil {
			return nil, err
		}

		endDate, err := time.Parse(time.RFC3339, conference.EndDate)
		if err != nil {
			return nil, err
		}

		var registrationDeadline *time.Time
		if conference.RegistrationDeadline != nil {
			parsedDeadline, err := time.Parse(time.RFC3339, *conference.RegistrationDeadline)
			if err != nil {
				return nil, err
			}
			registrationDeadline = &parsedDeadline
		}

		conferences = append(conferences, &models.Conference{
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
		})
	}

	return &models.ConferencesPage{
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

func GetConference(ctx context.Context, dbClient *db.DB, id int) (*models.Conference, error) {
	var c repositories.Conference
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		c, err = repositories.GetConference(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return &models.Conference{}, err
	}

	startDate, err := time.Parse(time.RFC3339, c.StartDate)
	if err != nil {
		return &models.Conference{}, err
	}

	endDate, err := time.Parse(time.RFC3339, c.EndDate)
	if err != nil {
		return &models.Conference{}, err
	}

	var registrationDeadline *time.Time
	if c.RegistrationDeadline != nil {
		parsedDeadline, err := time.Parse(time.RFC3339, *c.RegistrationDeadline)
		if err != nil {
			return &models.Conference{}, err
		}
		registrationDeadline = &parsedDeadline
	}

	return &models.Conference{
		ID:                   c.Id,
		Title:                c.Title,
		StartDate:            startDate,
		EndDate:              endDate,
		Location:             c.Location,
		Website:              c.Website,
		Acronym:              c.Acronym,
		AdditionalInfo:       c.AdditionalInfo,
		ParticipantsLimit:    c.ParticipantsLimit,
		RegistrationDeadline: registrationDeadline,
	}, nil
}

func GetConferencesMetrics(ctx context.Context, dbClient *db.DB) (*models.ConferencesMetrics, error) {
	var metrics repositories.ConferencesMetrics
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		metrics, err = repositories.GetMetrics(tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &models.ConferencesMetrics{
		RunningConferences:        metrics.RunningConferences,
		StartingInLessThan24Hours: metrics.StartingInLessThan24Hours,
		TotalConducted:            metrics.TotalConducted,
		ParticipantsToday:         metrics.ParticipantsToday,
	}, nil
}
