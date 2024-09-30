package repositories

import (
	"github.com/jmoiron/sqlx"
)

type AgendaItem struct {
	Id           int    `json:"id" db:"id"`
	ConferenceId int    `json:"conference_id" db:"conference_id"`
	StartTime    string `json:"start_time" db:"start_time"`
	EndTime      string `json:"end_time" db:"end_time"`
	Event        string `json:"event" db:"event"`
	Speaker      string `json:"speaker" db:"speaker"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

func (a *AgendaItem) TableName() string {
	return "public.agenda"
}

func GetAgenda(tx *sqlx.Tx, conferenceId int) ([]AgendaItem, error) {
	var agenda []AgendaItem
	a := AgendaItem{}
	query := "SELECT id, conference_id, start_time, end_time, event, speaker FROM " + a.TableName() + " WHERE conference_id = $1"
	err := tx.Select(
		&agenda,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return agenda, nil
}

func CreateAgenda(tx *sqlx.Tx, agenda AgendaItem) error {
	query := "INSERT INTO " + agenda.TableName() + " (conference_id, start_time, end_time, event, speaker) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(query, agenda.ConferenceId, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAgenda(tx *sqlx.Tx, agendaId int) error {
	query := "DELETE FROM " + new(AgendaItem).TableName() + " WHERE id = $1"
	_, err := tx.Exec(query, agendaId)
	if err != nil {
		return err
	}

	return nil
}

func CountAgendaItems(tx *sqlx.Tx, conferenceId int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + new(AgendaItem).TableName() + " WHERE conference_id = $1"
	err := tx.Get(&count, query, conferenceId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
