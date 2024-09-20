package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
)

type AgendaItem struct {
	Id           int    `json:"id" db:"id"`
	ConferenceId int    `json:"confrence_id" db:"conference_id"`
	StartTime    string `json:"start_time" db:"start_time"`
	EndTime      string `json:"end_time" db:"end_time"`
	Event        string `json:"event" db:"event"`
	Speaker      string `json:"speaker" db:"speaker"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

func (a *AgendaItem) TableName() string {
	return "agenda"
}

func GetAgenda(qe *db.QueryExecutor, conferenceId int) ([]AgendaItem, error) {
	var agenda []AgendaItem
	a := AgendaItem{}
	query := "SELECT id, conference_id, start_time, end_time, event, speaker FROM " + a.TableName() + " WHERE conference_id = ?"
	err := sqlx.Select(
		qe,
		&agenda,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return agenda, nil
}

func CreateAgenda(qe *db.QueryExecutor, agenda AgendaItem) error {
	query := "INSERT INTO " + agenda.TableName() + " (conference_id, start_time, end_time, event, speaker) VALUES (?, ?, ?, ?, ?)"
	_, err := qe.Exec(query, agenda.ConferenceId, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAgenda(qe *db.QueryExecutor, agenda AgendaItem) error {
	query := "UPDATE " + agenda.TableName() + " SET start_time = ?, end_time = ?, event = ?, speaker = ? WHERE id = ?"
	_, err := qe.Exec(query, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker, agenda.Id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAgenda(qe *db.QueryExecutor, agendaId int) error {
	a := AgendaItem{}
	query := "DELETE FROM " + a.TableName() + " WHERE id = ?"
	_, err := qe.Exec(query, agendaId)
	if err != nil {
		return err
	}

	return nil
}

func CountAgendaItems(qe *db.QueryExecutor, conferenceId int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + new(AgendaItem).TableName() + " WHERE conference_id = ?"
	err := sqlx.Get(qe, &count, query, conferenceId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
