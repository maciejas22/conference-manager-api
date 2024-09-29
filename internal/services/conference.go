package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	filters "github.com/maciejas22/conference-manager/api/db/repositories/shared"
	"github.com/maciejas22/conference-manager/api/internal/converters"
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
			file, err := converters.ConvertUploadFileToReader(f.UploadFile)
			if err != nil {
				return err
			}
			err = files.UploadConferenceFile(ctx, s3, conferenceId, file.Name, file.Content)
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
				file, err := converters.ConvertUploadFileToReader(f.UploadFile)
				if err != nil {
					return err
				}
				err = files.UploadConferenceFile(ctx, s3, input.ID, file.Name, file.Content)
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
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		c, m, err = repositories.GetAllConferences(
			tx,
			userId,
			converters.ConvertPageSchemaToRepo(p),
			converters.ConvertSortSchemaToRepo(s),
			converters.ConvertConferenceFiltersSchemaToRepo(f),
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
		conferences = append(conferences, converters.ConvertConferenceRepoToSchema(&conference))
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
	var conference repositories.Conference
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		conference, err = repositories.GetConference(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return &models.Conference{}, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
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
