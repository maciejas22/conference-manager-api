package models

import (
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/utils"
)

type AgendaItem struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime" db:"start_time"`
	EndTime   time.Time `json:"endTime"   db:"end_time"`
	Event     string    `json:"event"`
	Speaker   string    `json:"speaker"`
}

func (a *AgendaItem) ToRepo() *repositories.AgendaItem {
	startTimeString := utils.TimeToString(&a.StartTime)
	endTimeString := utils.TimeToString(&a.EndTime)
	return &repositories.AgendaItem{
		Id:        a.ID,
		StartTime: *startTimeString,
		EndTime:   *endTimeString,
		Event:     a.Event,
		Speaker:   a.Speaker,
	}
}
