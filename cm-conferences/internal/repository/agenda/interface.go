package repository

type AgendaRepoInterface interface {
	GetAgenda(conferenceId int) ([]AgendaItem, error)
	CreateAgendas(agenda []AgendaItem) ([]int, error)
	DeleteAgendas(agendaId []int) error
	CountAgendaItems(conferenceId []int) ([]EventsCount, error)
}
