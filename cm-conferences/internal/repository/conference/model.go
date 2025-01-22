package repository

import "time"

type Conference struct {
	Id                   int        `json:"id" db:"id"`
	Title                string     `json:"title" db:"title"`
	StartDate            time.Time  `json:"start_date" db:"start_date"`
	EndDate              time.Time  `json:"end_date" db:"end_date"`
	Location             string     `json:"location" db:"location"`
	Website              *string    `json:"website,omitempty" db:"website"`
	Acronym              *string    `json:"acronym,omitempty" db:"acronym"`
	AdditionalInfo       *string    `json:"additional_info,omitempty" db:"additional_info"`
	ParticipantsLimit    *int       `json:"participants_limit,omitempty" db:"participants_limit"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty" db:"registration_deadline"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	TicketPrice          int        `json:"ticket_price,omitempty" db:"ticket_price"`
}
