package repository

type ConferenceOrganizerRepoInterface interface {
	IsConferenceOrganizer(conferenceId int, userId int) (bool, error)
	AddConferenceOrganizer(conferenceId int, userId int) (bool, error)
	GetOrganizerLevelMetrics(organizerId int) (ConferenceOrganizerMetrics, error)
	GetParticipantsTrend(organizerId int) ([]TrendEntry, error)
	GetConferenceOrganizerId(conferenceId int) (int, error)
}
