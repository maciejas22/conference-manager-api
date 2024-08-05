package models

import (
	"time"
)

type AgendaItem struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime" db:"start_time"`
	EndTime   time.Time `json:"endTime"   db:"end_time"`
	Event     string    `json:"event"`
	Speaker   string    `json:"speaker"`
}
