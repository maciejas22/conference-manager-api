package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type ConferenceOrganizer struct {
	UserId       int    `json:"user_id" db:"user_id"`
	ConferenceId int    `json:"conference_id" db:"conference_id"`
	JoinedAt     string `json:"joined_at" db:"joined_at"`
}

func (c *ConferenceOrganizer) TableName() string {
	return "public.conference_organizers"
}

func IsConferenceOrganizer(tx *sqlx.Tx, conferenceId int, userId int) (bool, error) {
	var count int
	o := &ConferenceOrganizer{}
	query := "SELECT COUNT(*) FROM " + o.TableName() + " WHERE conference_id = $1 AND user_id = $2"
	err := tx.Get(
		&count,
		query,
		conferenceId,
		userId,
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func AddConferenceOrganizer(tx *sqlx.Tx, conferenceId int, userId int) (bool, error) {
	o := &ConferenceOrganizer{
		UserId:       userId,
		ConferenceId: conferenceId,
	}
	query := "INSERT INTO " + o.TableName() + " (user_id, conference_id) VALUES ($1, $2)"
	_, err := tx.Exec(
		query,
		o.UserId,
		o.ConferenceId,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

type ConferenceOrganizerMetrics struct {
	RunningConferencesCount   int     `db:"running_conferences_count"`
	ParticipantsCount         int     `db:"participants_count"`
	AverageParticipantsCount  float64 `db:"average_participants_count"`
	TotalOrganizedConferences int     `db:"total_organized_conferences"`
}

func GetOrganizerLevelMetrics(tx *sqlx.Tx, organizerId int) (ConferenceOrganizerMetrics, error) {
	var metrics ConferenceOrganizerMetrics
	o := &ConferenceOrganizer{}
	c := &Conference{}
	cp := &ConferenceParticipant{}
	query := `SELECT 
					(SELECT COUNT(*) FROM ` + c.TableName() + ` WHERE id IN (SELECT conference_id FROM ` + o.TableName() + ` WHERE user_id = $1) AND NOW() BETWEEN start_date AND end_date) AS running_conferences_count,
					(SELECT COUNT(*) FROM ` + cp.TableName() + ` WHERE conference_id IN (SELECT id FROM ` + c.TableName() + ` WHERE id IN (SELECT conference_id FROM ` + o.TableName() + ` WHERE user_id = $2))) AS participants_count,
					(SELECT COUNT(*) FROM ` + o.TableName() + ` WHERE user_id = $3) AS total_organized_conferences`
	err := tx.Get(
		&metrics,
		query,
		organizerId,
		organizerId,
		organizerId,
	)
	if err != nil {
		return ConferenceOrganizerMetrics{}, err
	}
	if metrics.RunningConferencesCount > 0 {
		metrics.AverageParticipantsCount = float64(metrics.ParticipantsCount) / float64(metrics.RunningConferencesCount)
	} else {
		metrics.AverageParticipantsCount = 0
	}
	return metrics, nil
}

type TrendEntry struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

func GetParticipantsTrend(tx *sqlx.Tx, organizerId int) ([]TrendEntry, error) {
	var counts []TrendEntry

	var reportStartDate *time.Time
	query := `
	  SELECT MIN(c.start_date)
	  FROM ` + (new(Conference)).TableName() + ` c
	  JOIN ` + (new(ConferenceOrganizer)).TableName() + ` o ON c.id = o.conference_id
	  WHERE o.user_id = $1 AND c.start_date <= NOW() AND c.end_date >= NOW()
  `
	err := tx.Get(&reportStartDate, query, organizerId)
	if err != nil {
		return counts, err
	}

	if reportStartDate == nil {
		return counts, nil
	}

	totalDuration := time.Since(*reportStartDate)
	interval := totalDuration / 10
	if interval.Hours() < 24 {
		interval = 24 * time.Hour
	}

	query = `
	  WITH intervals AS (
	  	SELECT generate_series(
	  		$1::timestamp,
        NOW(),
	  		$2::interval
	  	) start_time
	  )
	  SELECT
	  	i.start_time AS date,
	  	COALESCE(COUNT(cp.conference_id), 0) AS value 
	  FROM intervals i
	  LEFT JOIN ` + (new(ConferenceParticipant)).TableName() + ` cp
	  	ON cp.joined_at >= i.start_time
	  	AND cp.joined_at < i.start_time + $2::interval
	  LEFT JOIN ` + (new(ConferenceOrganizer)).TableName() + ` o
	  	ON cp.conference_id = o.conference_id
	  WHERE o.user_id = $3 OR cp.conference_id IS NULL
	  GROUP BY i.start_time
	  ORDER BY i.start_time ASC
  `

	err = tx.Select(&counts, query, reportStartDate, interval.String(), organizerId)
	if err != nil {
		return counts, err
	}

	return counts, nil
}

func GetConferenceOrganizerId(tx *sqlx.Tx, conferenceId int) (int, error) {
	var organizerId int
	o := &ConferenceOrganizer{}
	query := "SELECT user_id FROM " + o.TableName() + " WHERE conference_id = $1"
	err := tx.Get(
		&organizerId,
		query,
		conferenceId,
	)
	if err != nil {
		return 0, err
	}
	return organizerId, nil
}
