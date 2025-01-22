package repository

import "time"

type AgendaItem struct {
	Id           int       `json:"id" db:"id"`
	ConferenceId int       `json:"conference_id" db:"conference_id"`
	StartTime    time.Time `json:"start_time" db:"start_time"`
	EndTime      time.Time `json:"end_time" db:"end_time"`
	Event        string    `json:"event" db:"event"`
	Speaker      string    `json:"speaker" db:"speaker"`
	CreatedAt    string    `json:"created_at" db:"created_at"`
	UpdatedAt    string    `json:"updated_at" db:"updated_at"`
}

type EventsCount struct {
	ConferenceId int `json:"conference_id" db:"conference_id"`
	Count        int `json:"count" db:"count"`
}
