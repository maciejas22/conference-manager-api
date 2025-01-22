package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	table "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"
)

type AgendaRepo struct {
	Db *sqlx.DB
}

func NewAgendaRepo(db *sqlx.DB) AgendaRepoInterface {
	return &AgendaRepo{Db: db}
}

func (r *AgendaRepo) GetAgenda(conferenceId int) ([]AgendaItem, error) {
	var agenda []AgendaItem
	query := "SELECT id, conference_id, start_time, end_time, event, speaker FROM " + table.GetTableName(table.AgendaTable) + " WHERE conference_id = $1"
	err := r.Db.Select(
		&agenda,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, errors.New("Agenda not found")
	}
	return agenda, nil
}

func (r *AgendaRepo) CreateAgendas(agendas []AgendaItem) ([]int, error) {
	var ids []int
	query := "INSERT INTO " + table.GetTableName(table.AgendaTable) + " (conference_id, start_time, end_time, event, speaker) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	tx, err := r.Db.Beginx()
	if err != nil {
		return nil, errors.New("Could not start agenda create transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, agenda := range agendas {
		var id int
		err := tx.QueryRowx(query, agenda.ConferenceId, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker).Scan(&id)
		if err != nil {
			return nil, errors.New("Could not insert agenda")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *AgendaRepo) DeleteAgendas(agendaIds []int) error {
	query := "DELETE FROM " + table.GetTableName(table.AgendaTable) + " WHERE id = $1"
	tx, err := r.Db.Beginx()
	if err != nil {
		return errors.New("Could not start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, id := range agendaIds {
		_, err := tx.Exec(query, id)
		if err != nil {
			return errors.New("Could not delete agenda")
		}
	}

	return nil
}

func (r *AgendaRepo) CountAgendaItems(conferenceIds []int) ([]EventsCount, error) {
	results := []EventsCount{}
	query := "SELECT conference_id, COUNT(*) as count FROM " + table.GetTableName(table.AgendaTable) + " WHERE conference_id = ANY($1) GROUP BY conference_id"

	err := r.Db.Select(&results, query, conferenceIds)
	if err != nil {
		return nil, errors.New("Could not get agenda items count")
	}
	return results, nil
}
