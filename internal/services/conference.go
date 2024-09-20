package services

import (
	"context"
	"strconv"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/internal/utils"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func CreateConference(ctx context.Context, dbClient *db.DB, s3 *s3.S3Client, userId int, createConferenceInput models.CreateConferenceInput) (*int, error) {
	var conferenceId int
	err := db.Transaction(ctx, dbClient.QueryExecutor, func(qe *db.QueryExecutor) error {
		startDate := createConferenceInput.StartDate
		startDateString := utils.TimeToString(&startDate)
		endDate := createConferenceInput.EndDate
		endDateString := utils.TimeToString(&endDate)
		deadlineString := utils.TimeToString(createConferenceInput.RegistrationDeadline)
		conferenceId, err := repositories.CreateConference(qe, repositories.Conference{
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
			startTime := utils.TimeToString(&a.StartTime)
			endTime := utils.TimeToString(&a.EndTime)
			err = repositories.CreateAgenda(qe, repositories.AgendaItem{
				ConferenceId: conferenceId,
				StartTime:    *startTime,
				EndTime:      *endTime,
				Event:        a.Event,
				Speaker:      a.Speaker,
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
			err = repositories.UploadFile(ctx, s3, strconv.Itoa(conferenceId), file.Name, file.Content)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &conferenceId, err
}

func ModifyConference(ctx context.Context, dbClient *db.DB, s3 *s3.S3Client, input models.ModifyConferenceInput) (*int, error) {
	var conferenceId *int
	err := db.Transaction(ctx, dbClient.QueryExecutor, func(qe *db.QueryExecutor) error {
		startDate := input.StartDate
		startDateString := utils.TimeToString(startDate)
		endDate := input.EndDate
		endDateString := utils.TimeToString(endDate)
		deadlineString := utils.TimeToString(input.RegistrationDeadline)
		conference, err := repositories.UpdateConference(qe, repositories.Conference{
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
			startTime := utils.TimeToString(&a.StartTime)
			endTime := utils.TimeToString(&a.EndTime)

			if a.ID != nil && a.Destroy != nil && *a.Destroy {
				err = repositories.DeleteAgenda(qe, *a.ID)
				if err != nil {
					return err
				}
			} else if a.ID == nil {
				err = repositories.CreateAgenda(qe, repositories.AgendaItem{
					ConferenceId: conference.Id,
					StartTime:    *startTime,
					EndTime:      *endTime,
					Event:        a.Event,
					Speaker:      a.Speaker,
				})

				if err != nil {
					return err
				}

			} else if a.ID != nil {
				err = repositories.UpdateAgenda(qe, repositories.AgendaItem{
					Id:           *a.ID,
					ConferenceId: conference.Id,
					StartTime:    *startTime,
					EndTime:      *endTime,
					Event:        a.Event,
					Speaker:      a.Speaker,
				})

				if err != nil {
					return err
				}

			}
		}

		for _, f := range input.Files {
			if f.DeleteFile != nil {
				err = repositories.DeleteFile(ctx, s3, strconv.Itoa(f.DeleteFile.ID))
				if err != nil {
					return err
				}
			} else if f.UploadFile != nil {
				file, err := converters.ConvertUploadFileToReader(f.UploadFile)
				if err != nil {
					return err
				}
				err = repositories.UploadFile(ctx, s3, strconv.Itoa(conference.Id), file.Name, file.Content)
				if err != nil {
					return err
				}
			}
		}

		conferenceId = &conference.Id
		return nil
	})

	return conferenceId, err
}

func GetAllConferences(ctx context.Context, dbClient *db.DB, userId int, page *models.Page, sort *models.Sort, filters *models.ConferenceFilter) (*models.ConferencePage, error) {
	c, m, err := repositories.GetAllConferences(
		dbClient.QueryExecutor,
		userId,
		converters.ConvertPageSchemaToRepo(page),
		converters.ConvertSortSchemaToRepo(sort),
		converters.ConvertConferenceFiltersSchemaToRepo(filters),
	)
	if err != nil {
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

func GetConference(ctx context.Context, dbClient *db.DB, id int) (*models.Conference, error) {
	conference, err := repositories.GetConference(dbClient.QueryExecutor, id)
	if err != nil {
		return &models.Conference{}, err
	}

	return converters.ConvertConferenceRepoToSchema(&conference), nil
}

func GetConferencesMetrics(ctx context.Context, dbClient *db.DB) (*models.ConferencesMetrics, error) {
	metrics, err := repositories.GetMetrics(dbClient.QueryExecutor)
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
