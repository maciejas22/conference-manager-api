package repository

import "time"

type ConferenceOrganizer struct {
	UserId       int       `json:"user_id" db:"user_id"`
	ConferenceId int       `json:"conference_id" db:"conference_id"`
	JoinedAt     time.Time `json:"joined_at" db:"joined_at"`
}
