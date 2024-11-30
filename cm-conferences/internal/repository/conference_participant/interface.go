package repository

import common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"

type ConferenceParticipantRepoInterface interface {
	GetConferenceParticipantsCount(conferenceIds []int) ([]ParticipantsCount, error)
	AddConferenceParticipant(conferenceId int, userId int) (string, error)
	RemoveConferenceParticipant(conferenceId int, userId int) (int, error)
	IsConferenceParticipant(conferenceId int, userId int) (bool, error)
	GetParticipantTickets(userId int, page common.Page) ([]ConferenceParticipant, common.PaginationMeta, error)
	IsTicketValid(ticketId string, conferenceId int) (bool, error)
}
