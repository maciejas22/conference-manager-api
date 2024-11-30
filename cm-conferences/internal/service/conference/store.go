package service

import (
	"context"

	agendaRepo "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/agenda"
	repoFilters "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"
	confRepo "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference"
	common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/common"
)

type ConferenceService struct {
	conferenceRepo confRepo.ConferenceRepoInterface
	agendaRepo     agendaRepo.AgendaRepoInterface
}

func NewConferenceService(conferenceRepo confRepo.ConferenceRepoInterface, agendaRepo agendaRepo.AgendaRepoInterface) ConferenceServiceInterface {
	return &ConferenceService{
		conferenceRepo: conferenceRepo,
		agendaRepo:     agendaRepo,
	}
}

func (s *ConferenceService) CreateConference(ctx context.Context, userId int, createConferenceInput CreateConferenceInput) (*int, error) {
	c, err := s.conferenceRepo.CreateConference(confRepo.Conference{
		Title:                createConferenceInput.Title,
		StartDate:            createConferenceInput.StartDate,
		EndDate:              createConferenceInput.EndDate,
		Location:             createConferenceInput.Location,
		Website:              createConferenceInput.Website,
		Acronym:              createConferenceInput.Acronym,
		AdditionalInfo:       createConferenceInput.AdditionalInfo,
		ParticipantsLimit:    createConferenceInput.ParticipantsLimit,
		RegistrationDeadline: createConferenceInput.RegistrationDeadline,
		TicketPrice:          createConferenceInput.TicketPrice,
	}, userId)
	if err != nil {
		return nil, err
	}

	agendaItems := make([]agendaRepo.AgendaItem, len(createConferenceInput.Agenda))
	for i, a := range createConferenceInput.Agenda {
		agendaItems[i] = agendaRepo.AgendaItem{
			ConferenceId: c,
			StartTime:    a.StartTime,
			EndTime:      a.EndTime,
			Event:        a.Event,
			Speaker:      a.Speaker,
		}
	}
	_, err = s.agendaRepo.CreateAgendas(agendaItems)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *ConferenceService) ModifyConference(ctx context.Context, input ModifyConferenceInput) (*int, error) {
	c, err := s.conferenceRepo.UpdateConference(confRepo.Conference{
		Id:                   input.ID,
		Title:                *input.Title,
		StartDate:            *input.StartDate,
		EndDate:              *input.EndDate,
		Location:             *input.Location,
		Website:              input.Website,
		Acronym:              input.Acronym,
		AdditionalInfo:       input.AdditionalInfo,
		ParticipantsLimit:    input.ParticipantsLimit,
		RegistrationDeadline: input.RegistrationDeadline,
		TicketPrice:          *input.TicketPrice,
	})
	if err != nil {
		return nil, err
	}

	var agendaItemsToCreate []agendaRepo.AgendaItem
	var agendaItemsToDelete []int
	for _, a := range input.Agenda {
		if a.Id != 0 {
			agendaItemsToDelete = append(agendaItemsToDelete, a.Id)
		} else {
			agendaItemsToCreate = append(agendaItemsToCreate, agendaRepo.AgendaItem{
				ConferenceId: c,
				StartTime:    a.StartTime,
				EndTime:      a.EndTime,
				Event:        a.Event,
				Speaker:      a.Speaker,
			})
		}
	}
	if len(agendaItemsToDelete) > 0 {
		err = s.agendaRepo.DeleteAgendas(agendaItemsToDelete)
		if err != nil {
			return nil, err
		}
	}
	if len(agendaItemsToCreate) > 0 {
		_, err = s.agendaRepo.CreateAgendas(agendaItemsToCreate)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (s *ConferenceService) GetConferencesPage(ctx context.Context, userId int, p *common.Page, sort *common.Sort, f *ConferencesFilters) ([]int, common.PaginationMeta, error) {
	var cPage *repoFilters.Page
	if p != nil {
		cPage = &repoFilters.Page{
			PageNumber: p.PageNumber,
			PageSize:   p.PageSize,
		}
	} else {
		cPage = &repoFilters.Page{
			PageNumber: 1,
			PageSize:   10,
		}
	}

	var cSort *repoFilters.Sort
	if sort != nil {
		var order repoFilters.Order
		if sort.Order == common.ASC {
			order = repoFilters.ASC
		} else {
			order = repoFilters.DESC
		}

		cSort = &repoFilters.Sort{
			Column: sort.Column,
			Order:  order,
		}
	} else {
		cSort = &repoFilters.Sort{
			Column: "id",
			Order:  repoFilters.ASC,
		}
	}

	var cFilters *confRepo.ConferenceFilter
	if f != nil {
		cFilters = &confRepo.ConferenceFilter{
			Title:          f.Title,
			AssociatedOnly: f.AssociatedOnly,
			RunningOnly:    f.RunningOnly,
		}
	} else {
		cFilters = &confRepo.ConferenceFilter{}
	}

	c, m, err := s.conferenceRepo.GetConferencesPage(userId, *cPage, cSort, cFilters)
	if err != nil {
		return nil, common.PaginationMeta{}, err
	}

	return c, common.PaginationMeta{
		PageNumber: m.PageNumber,
		PageSize:   m.PageSize,
		TotalItems: m.TotalItems,
		TotalPages: m.TotalPages,
	}, nil
}

func (s *ConferenceService) GetConference(ctx context.Context, id int) (Conference, error) {
	c, err := s.conferenceRepo.GetConferenceById(id)
	if err != nil {
		return Conference{}, err
	}

	return Conference{
		Id:                   c.Id,
		Title:                c.Title,
		StartDate:            c.StartDate,
		EndDate:              c.EndDate,
		Location:             c.Location,
		Website:              c.Website,
		Acronym:              c.Acronym,
		AdditionalInfo:       c.AdditionalInfo,
		ParticipantsLimit:    c.ParticipantsLimit,
		RegistrationDeadline: c.RegistrationDeadline,
		TicketPrice:          c.TicketPrice,
	}, nil
}

func (s *ConferenceService) GetConferencesMetrics(ctx context.Context) (ConferencesMetrics, error) {
	m, err := s.conferenceRepo.GetMetrics()
	if err != nil {
		return ConferencesMetrics{}, err
	}

	return ConferencesMetrics{
		RunningConferences:        m.RunningConferences,
		StartingInLessThan24Hours: m.StartingInLessThan24Hours,
		TotalConducted:            m.TotalConducted,
		ParticipantsToday:         m.ParticipantsToday,
	}, nil
}

func (s *ConferenceService) GetAgenda(ctx context.Context, conferenceId int) ([]AgendaItem, error) {
	a, err := s.agendaRepo.GetAgenda(conferenceId)
	if err != nil {
		return nil, err
	}

	agenda := make([]AgendaItem, len(a))
	for i, item := range a {
		agenda[i] = AgendaItem{
			Id:           item.Id,
			ConferenceId: item.ConferenceId,
			StartTime:    item.StartTime,
			EndTime:      item.EndTime,
			Event:        item.Event,
			Speaker:      item.Speaker,
		}
	}

	return agenda, nil
}

func (s *ConferenceService) GetConferencesByIds(ctx context.Context, ids []int) ([]Conference, error) {
	c, err := s.conferenceRepo.GetConferencesByIds(ids)
	if err != nil {
		return nil, err
	}

	conferences := make([]Conference, len(c))
	for i, item := range c {
		conferences[i] = Conference{
			Id:                   item.Id,
			Title:                item.Title,
			StartDate:            item.StartDate,
			EndDate:              item.EndDate,
			Location:             item.Location,
			Website:              item.Website,
			Acronym:              item.Acronym,
			AdditionalInfo:       item.AdditionalInfo,
			ParticipantsLimit:    item.ParticipantsLimit,
			RegistrationDeadline: item.RegistrationDeadline,
			TicketPrice:          item.TicketPrice,
		}
	}

	return conferences, nil
}
