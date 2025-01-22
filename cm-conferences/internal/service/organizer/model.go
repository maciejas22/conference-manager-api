package service

import "time"

type ConferenceOrganizer struct {
	UserId       int
	ConferenceId int
}

type OrganizerMetrics struct {
	RunningConferences        int
	ParticipantsCount         int
	AverageParticipantsCount  float64
	TotalOrganizedConferences int
}

type ParticipantsTrendEntry struct {
	Date            time.Time
	NewParticipants int
}

type NewParticipantsTrend struct {
	Entries []ParticipantsTrendEntry
}
