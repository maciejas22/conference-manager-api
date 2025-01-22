package repository

import "time"

type ConferenceParticipant struct {
	UserId       int       `json:"user_id" db:"user_id"`
	ConferenceId int       `json:"conference_id" db:"conference_id"`
	TicketId     string    `json:"ticket_id" db:"ticket_id"`
	JoinedAt     time.Time `json:"joined_at" db:"joined_at"`
}

type ParticipantsCount struct {
	ConferenceId int `db:"conference_id"`
	Count        int `db:"count"`
}
