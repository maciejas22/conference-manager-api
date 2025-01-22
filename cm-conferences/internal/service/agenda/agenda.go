package service

import (
	repo "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/agenda"
)

type AgendaService struct {
	agendaRepo repo.AgendaRepoInterface
}

func NewAgendaService(agendaRepo repo.AgendaRepoInterface) AgendaServiceInterface {
	return &AgendaService{
		agendaRepo: agendaRepo,
	}
}

func (s *AgendaService) GetAgenda(conferenceId int) ([]AgendaItem, error) {
	agenda, err := s.agendaRepo.GetAgenda(conferenceId)
	if err != nil {
		return nil, err
	}

	var agendaItems []AgendaItem
	for _, a := range agenda {
		agendaItems = append(agendaItems, AgendaItem{
			Id:        a.Id,
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
			Event:     a.Event,
			Speaker:   a.Speaker,
		})
	}

	return agendaItems, nil
}

func (s *AgendaService) GetAgendaItemsCount(conferenceIds []int) ([]AgendaItemsCount, error) {
	counts, err := s.agendaRepo.CountAgendaItems(conferenceIds)
	if err != nil {
		return nil, err
	}

	countsAgendaItems := make([]AgendaItemsCount, len(counts))
	for i, c := range counts {
		countsAgendaItems[i] = AgendaItemsCount{
			ConferenceId:     c.ConferenceId,
			AgendaItemsCount: c.Count,
		}
	}
	return countsAgendaItems, nil
}
