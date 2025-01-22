package service

import "time"

type AgendaItem struct {
	Id           int
	ConferenceId int
	StartTime    time.Time
	EndTime      time.Time
	Event        string
	Speaker      string
}

type AgendaItemsCount struct {
	ConferenceId     int
	AgendaItemsCount int
}
