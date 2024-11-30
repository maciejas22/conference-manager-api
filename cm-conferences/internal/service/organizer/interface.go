package service

type OrganizerServiceInterface interface {
	GetConferenceOrganizerId(conferenceId int) (userId int, err error)
	GetOrganizerMetrics(organizerId int) (*OrganizerMetrics, error)
	GetParticipantsJoiningTrend(organizerId int) (*NewParticipantsTrend, error)
	IsConferenceOrganizer(userId, conferenceID int) (bool, error)
	AddOrganizerToConference(userId int, conferenceID int) (bool, error)
}
