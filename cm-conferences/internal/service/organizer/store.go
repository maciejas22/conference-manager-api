package service

import (
	repo "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference_organizer"
)

type OrganizerService struct {
	organizerRepo repo.ConferenceOrganizerRepoInterface
}

func NewConferenceOrganizerService(conferenceOrganizerRepo repo.ConferenceOrganizerRepoInterface) OrganizerServiceInterface {
	return &OrganizerService{conferenceOrganizerRepo}
}

func (s *OrganizerService) GetConferenceOrganizerId(conferenceId int) (int, error) {
	id, err := s.organizerRepo.GetConferenceOrganizerId(conferenceId)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *OrganizerService) GetOrganizerMetrics(organizerId int) (*OrganizerMetrics, error) {
	organizerMetrics, err := s.organizerRepo.GetOrganizerLevelMetrics(organizerId)
	if err != nil {
		return nil, err
	}

	return &OrganizerMetrics{
		RunningConferences:        organizerMetrics.RunningConferencesCount,
		ParticipantsCount:         organizerMetrics.ParticipantsCount,
		AverageParticipantsCount:  organizerMetrics.AverageParticipantsCount,
		TotalOrganizedConferences: organizerMetrics.TotalOrganizedConferences,
	}, nil
}

func (s *OrganizerService) GetParticipantsJoiningTrend(organizerId int) (*NewParticipantsTrend, error) {
	participantsJoiningTrend, err := s.organizerRepo.GetParticipantsTrend(organizerId)
	if err != nil {
		return nil, err
	}

	var trendEntries NewParticipantsTrend
	for _, trendEntry := range participantsJoiningTrend {
		trendEntries.Entries = append(trendEntries.Entries, ParticipantsTrendEntry{
			Date:            trendEntry.Date,
			NewParticipants: trendEntry.Value,
		})
	}

	return &trendEntries, nil
}

func (s *OrganizerService) IsConferenceOrganizer(userId, conferenceID int) (bool, error) {
	isOrganizer, err := s.organizerRepo.IsConferenceOrganizer(conferenceID, userId)
	if err != nil {
		return false, err
	}

	return isOrganizer, nil
}

func (s *OrganizerService) AddOrganizerToConference(userId int, conferenceID int) (bool, error) {
	oId, err := s.organizerRepo.AddConferenceOrganizer(conferenceID, userId)
	if err != nil {
		return false, err
	}

	return oId, nil
}
