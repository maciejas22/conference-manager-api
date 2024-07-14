package models

import (
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/utils"
)

type Conference struct {
	ID                   string     `json:"id"`
	Title                string     `json:"title"`
	Date                 time.Time  `json:"date"`
	Location             string     `json:"location"`
	AdditionalInfo       *string    `json:"additional_info,omitempty"       db:"additional_info"`
	ParticipantsLimit    *int       `json:"participants_limit,omitempty"    db:"participants_limit"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty" db:"registration_deadline"`
}

func (c *Conference) ToRepo() *repositories.Conference {
	dateString := utils.TimeToString(&c.Date)
	deadlineString := utils.TimeToString(c.RegistrationDeadline)
	return &repositories.Conference{
		Id:                   c.ID,
		Title:                c.Title,
		Date:                 *dateString,
		Location:             c.Location,
		AdditionalInfo:       c.AdditionalInfo,
		ParticipantsLimit:    c.ParticipantsLimit,
		RegistrationDeadline: deadlineString,
	}
}

type ConferenceFilter struct {
	AssociatedOnly *bool   `json:"associatedOnly,omitempty"`
	Title          *string `json:"title,omitempty"`
}

func (f *ConferenceFilter) ToRepo() *repositories.ConferenceFilter {
	if f == nil {
		return nil
	}

	return &repositories.ConferenceFilter{
		AssociatedOnly: f.AssociatedOnly,
		Title:          f.Title,
	}
}
