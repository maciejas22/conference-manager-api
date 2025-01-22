package service

import p "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/common"

type ParticipantServiceInterface interface {
	GetParticipantsCount(conferenceIds []int) ([]ConferenceParticipantsCount, error)
	AddUserToConference(userId int, conferenceID int) (string, error)
	RemoveUserFromConference(userId int, conferenceID int) (int, error)
	IsConferenceParticipant(userId, conferenceID int) (bool, error)
	GetParticipantsTickets(participantId int, page p.Page) ([]Ticket, p.PaginationMeta, error)
	IsTicketValid(ticketId string, conferenceID int) (bool, error)
}
