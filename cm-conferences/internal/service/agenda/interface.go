package service

type AgendaServiceInterface interface {
	GetAgenda(conferenceId int) ([]AgendaItem, error)
	GetAgendaItemsCount(conferenceIds []int) ([]AgendaItemsCount, error)
}
