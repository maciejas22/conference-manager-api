package service

import (
	repoFilters "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"
	repo "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference_participant"
	p "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/common"
)

type ParticipantService struct {
	participantRepo repo.ConferenceParticipantRepoInterface
}

func NewParticipantService(participantRepo repo.ConferenceParticipantRepoInterface) ParticipantServiceInterface {
	return &ParticipantService{participantRepo}
}

func (s *ParticipantService) GetParticipantsCount(conferenceIds []int) ([]ConferenceParticipantsCount, error) {
	participantsCount, err := s.participantRepo.GetConferenceParticipantsCount(conferenceIds)
	if err != nil {
		return nil, err
	}

	pc := make([]ConferenceParticipantsCount, len(participantsCount))
	for i, p := range participantsCount {
		pc[i] = ConferenceParticipantsCount{
			ConferenceId:      p.ConferenceId,
			ParticipantsCount: p.Count,
		}
	}
	return pc, nil
}

func (s *ParticipantService) AddUserToConference(userId int, conferenceID int) (string, error) {
	tId, err := s.participantRepo.AddConferenceParticipant(conferenceID, userId)
	if err != nil {
		return "", err
	}

	return tId, nil
}

func (s *ParticipantService) RemoveUserFromConference(userId int, conferenceID int) (int, error) {
	cId, err := s.participantRepo.RemoveConferenceParticipant(conferenceID, userId)
	if err != nil {
		return 0, err
	}

	return cId, nil
}

func (s *ParticipantService) IsConferenceParticipant(userId, conferenceID int) (bool, error) {
	isParticipant, err := s.participantRepo.IsConferenceParticipant(conferenceID, userId)
	if err != nil {
		return false, err
	}

	return isParticipant, nil
}

func (s *ParticipantService) GetParticipantsTickets(participantId int, page p.Page) ([]Ticket, p.PaginationMeta, error) {
	ticketPage := repoFilters.Page{
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
	}
	ticketData, ticketsMeta, err := s.participantRepo.GetParticipantTickets(participantId, ticketPage)
	if err != nil {
		return nil, p.PaginationMeta{}, err
	}

	tickets := make([]Ticket, len(ticketData))
	for i, ticket := range ticketData {
		tickets[i] = Ticket{
			Id:           ticket.TicketId,
			ConferenceId: ticket.ConferenceId,
		}
	}
	meta := p.PaginationMeta{
		PageNumber: ticketsMeta.PageNumber,
		PageSize:   ticketsMeta.PageSize,
		TotalItems: ticketsMeta.TotalItems,
		TotalPages: ticketsMeta.TotalPages,
	}

	return tickets, meta, nil
}

func (s *ParticipantService) IsTicketValid(ticketId string, conferenceID int) (bool, error) {
	isValid, err := s.participantRepo.IsTicketValid(ticketId, conferenceID)
	if err != nil {
		return false, err
	}

	return isValid, nil
}
