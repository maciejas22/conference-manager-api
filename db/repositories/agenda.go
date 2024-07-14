package repositories

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
)

type AgendaItem struct {
	Id           string `json:"id" db:"id"`
	ConferenceId string `json:"confrence_id" db:"conference_id"`
	StartTime    string `json:"start_time" db:"start_time"`
	EndTime      string `json:"end_time" db:"end_time"`
	Event        string `json:"event" db:"event"`
	Speaker      string `json:"speaker" db:"speaker"`
}

func (a *AgendaItem) TableName() string {
	return "public.agenda"
}

type AgendaRepository interface {
	GetAgenda(conferenceId string) ([]AgendaItem, error)
	CreateAgenda(agenda AgendaItem) (AgendaItem, error)
	UpdateAgenda(agenda AgendaItem) (AgendaItem, error)
	DeleteAgenda(agendaId string) error
}

type agendaRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewAgendaRepository(ctx context.Context, db *db.DB) AgendaRepository {
	return &agendaRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *agendaRepository) GetAgenda(conferenceId string) ([]AgendaItem, error) {
	var agenda []AgendaItem
	a := AgendaItem{}
	query := "SELECT id, conference_id, start_time, end_time, event, speaker FROM " + a.TableName() + " WHERE conference_id = $1"
	err := r.db.SqlConn.Select(
		&agenda,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return agenda, nil
}

func (r *agendaRepository) CreateAgenda(agenda AgendaItem) (AgendaItem, error) {
	a := AgendaItem{}
	query := "INSERT INTO " + a.TableName() + " (conference_id, start_time, end_time, event, speaker) VALUES ($1, $2, $3, $4, $5) RETURNING id, conference_id, start_time, end_time, event, speaker"
	err := r.db.SqlConn.QueryRowx(query, agenda.ConferenceId, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker).Scan(&a.Id, &a.ConferenceId, &a.StartTime, &a.EndTime, &a.Event, &a.Speaker)
	if err != nil {
		return AgendaItem{}, err
	}

	return a, nil
}

func (r *agendaRepository) UpdateAgenda(agenda AgendaItem) (AgendaItem, error) {
	a := AgendaItem{}
	query := "UPDATE " + a.TableName() + " SET start_time = $1, end_time = $2, event = $3, speaker = $4 WHERE id = $5 RETURNING id, conference_id, start_time, end_time, event, speaker"
	err := r.db.SqlConn.QueryRowx(query, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker, agenda.Id).Scan(&a.Id, &a.ConferenceId, &a.StartTime, &a.EndTime, &a.Event, &a.Speaker)
	if err != nil {
		return AgendaItem{}, err
	}

	return a, nil
}

func (r *agendaRepository) DeleteAgenda(agendaId string) error {
	a := AgendaItem{}
	query := "DELETE FROM " + a.TableName() + " WHERE id = $1"
	_, err := r.db.SqlConn.Exec(query, agendaId)
	if err != nil {
		return err
	}

	return nil
}
