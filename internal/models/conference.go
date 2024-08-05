package models

import (
	"time"
)

type Conference struct {
	ID                   string     `json:"id"`
	Title                string     `json:"title"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              time.Time  `json:"end_date"`
	Location             string     `json:"location"`
	Website              *string    `json:"website,omitempty"              db:"website"`
	Acronym              *string    `json:"acronym,omitempty"              db:"acronym"`
	AdditionalInfo       *string    `json:"additional_info,omitempty"       db:"additional_info"`
	ParticipantsLimit    *int       `json:"participants_limit,omitempty"    db:"participants_limit"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty" db:"registration_deadline"`
	Files                []*File    `json:"files,omitempty"`
}

type ConferenceFilter struct {
	AssociatedOnly *bool   `json:"associatedOnly,omitempty"`
	Title          *string `json:"title,omitempty"`
}
