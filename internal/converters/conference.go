package converters

import (
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertConferenceSchemaToRepo(c *models.Conference) *repositories.Conference {
	var registrationDeadline *string
	if c.RegistrationDeadline != nil {
		formattedDeadline := c.RegistrationDeadline.Format(time.RFC3339)
		registrationDeadline = &formattedDeadline
	}

	return &repositories.Conference{
		Id:                   c.ID,
		Title:                c.Title,
		StartDate:            c.StartDate.Format(time.RFC3339),
		EndDate:              c.EndDate.Format(time.RFC3339),
		Location:             c.Location,
		Website:              c.Website,
		Acronym:              c.Acronym,
		AdditionalInfo:       c.AdditionalInfo,
		ParticipantsLimit:    c.ParticipantsLimit,
		RegistrationDeadline: registrationDeadline,
	}
}

func ConvertConferenceRepoToSchema(c *repositories.Conference) *models.Conference {
	startDate, err := time.Parse(time.RFC3339, c.StartDate)
	if err != nil {
		return &models.Conference{}
	}

	endDate, err := time.Parse(time.RFC3339, c.EndDate)
	if err != nil {
		return &models.Conference{}
	}

	var registrationDeadline *time.Time
	if c.RegistrationDeadline != nil {
		parsedDeadline, err := time.Parse(time.RFC3339, *c.RegistrationDeadline)
		if err != nil {
			return &models.Conference{}
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
		Files:                nil,
	}
}

func ConvertConferenceFiltersSchemaToRepo(f *models.ConferenceFilter) *repositories.ConferenceFilter {
	if f == nil {
		return nil
	}

	return &repositories.ConferenceFilter{
		AssociatedOnly: f.AssociatedOnly,
		Title:          f.Title,
	}
}
