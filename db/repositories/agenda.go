package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
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

func GetAgenda(tx *sqlx.Tx, conferenceId string) ([]AgendaItem, error) {
	var agenda []AgendaItem
	a := AgendaItem{}
	query := "SELECT id, conference_id, start_time, end_time, event, speaker FROM " + a.TableName() + " WHERE conference_id = $1"
	err := tx.Select(
		&agenda,
		query,
		conferenceId,
	)
	if err != nil {
		log.Panicln("Error while fetching agenda items: ", err)
		return nil, err
	}
	return agenda, nil
}

func CreateAgenda(tx *sqlx.Tx, agenda AgendaItem) (AgendaItem, error) {
	a := AgendaItem{}
	query := "INSERT INTO " + a.TableName() + " (conference_id, start_time, end_time, event, speaker) VALUES ($1, $2, $3, $4, $5) RETURNING id, conference_id, start_time, end_time, event, speaker"
	err := tx.Get(&a, query, agenda.ConferenceId, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker)
	if err != nil {
		log.Panicln("Error while creating agenda item: ", err)
		return AgendaItem{}, err
	}

	return a, nil
}

func UpdateAgenda(tx *sqlx.Tx, agenda AgendaItem) (AgendaItem, error) {
	a := AgendaItem{}
	query := "UPDATE " + a.TableName() + " SET start_time = $1, end_time = $2, event = $3, speaker = $4 WHERE id = $5 RETURNING id, conference_id, start_time, end_time, event, speaker"
	err := tx.Get(&a, query, agenda.StartTime, agenda.EndTime, agenda.Event, agenda.Speaker, agenda.Id)
	if err != nil {
		log.Panicln("Error while updating agenda item: ", err)
		return AgendaItem{}, err
	}

	return a, nil
}

func DeleteAgenda(tx *sqlx.Tx, agendaId string) error {
	a := AgendaItem{}
	query := "DELETE FROM " + a.TableName() + " WHERE id = $1"
	_, err := tx.Exec(query, agendaId)
	if err != nil {
		log.Panicln("Error while deleting agenda item: ", err)
		return err
	}

	return nil
}

func CountAgendaItems(tx *sqlx.Tx, conferenceId string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + new(AgendaItem).TableName() + " WHERE conference_id = $1"
	err := tx.Get(&count, query, conferenceId)
	if err != nil {
		log.Panicln("Error while counting agenda items: ", err)
		return 0, err
	}
	return count, nil
}
