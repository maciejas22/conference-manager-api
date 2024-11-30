package repository

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"
)

type ConferenceOrganizerRepo struct {
	Db *sqlx.DB
}

func NewOrganizerRepo(db *sqlx.DB) ConferenceOrganizerRepoInterface {
	return &ConferenceOrganizerRepo{Db: db}
}

func (r *ConferenceOrganizerRepo) IsConferenceOrganizer(conferenceId int, userId int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + common.GetTableName(common.ConferenceOrganizerTable) + " WHERE conference_id = $1 AND user_id = $2"
	err := r.Db.Get(
		&count,
		query,
		conferenceId,
		userId,
	)
	if err != nil {
		return false, errors.New("Could not get organizer count")
	}
	return count > 0, nil
}

func (r *ConferenceOrganizerRepo) AddConferenceOrganizer(conferenceId int, userId int) (bool, error) {
	o := &ConferenceOrganizer{
		UserId:       userId,
		ConferenceId: conferenceId,
	}
	query := "INSERT INTO " + common.GetTableName(common.ConferenceOrganizerTable) + " (user_id, conference_id) VALUES ($1, $2)"
	_, err := r.Db.Exec(
		query,
		o.UserId,
		o.ConferenceId,
	)
	if err != nil {
		return false, errors.New("Could not insert organizer")
	}
	return true, nil
}

type ConferenceOrganizerMetrics struct {
	RunningConferencesCount   int     `db:"running_conferences_count"`
	ParticipantsCount         int     `db:"participants_count"`
	AverageParticipantsCount  float64 `db:"average_participants_count"`
	TotalOrganizedConferences int     `db:"total_organized_conferences"`
}

func (r *ConferenceOrganizerRepo) GetOrganizerLevelMetrics(organizerId int) (ConferenceOrganizerMetrics, error) {
	var metrics ConferenceOrganizerMetrics
	query := `SELECT 
					(SELECT COUNT(*) FROM ` + common.GetTableName(common.ConferenceTable) + ` WHERE id IN (SELECT conference_id FROM ` + common.GetTableName(common.ConferenceOrganizerTable) + ` WHERE user_id = $1) AND NOW() BETWEEN start_date AND end_date) AS running_conferences_count,
					(SELECT COUNT(*) FROM ` + common.GetTableName(common.ConferenceParticipantTable) + ` WHERE conference_id IN (SELECT id FROM ` + common.GetTableName(common.ConferenceTable) + ` WHERE id IN (SELECT conference_id FROM ` + common.GetTableName(common.ConferenceOrganizerTable) + ` WHERE user_id = $2))) AS participants_count,
					(SELECT COUNT(*) FROM ` + common.GetTableName(common.ConferenceOrganizerTable) + ` WHERE user_id = $3) AS total_organized_conferences`
	err := r.Db.Get(
		&metrics,
		query,
		organizerId,
		organizerId,
		organizerId,
	)
	if err != nil {
		return ConferenceOrganizerMetrics{}, errors.New("Could not get organizer metrics")
	}
	if metrics.RunningConferencesCount > 0 {
		metrics.AverageParticipantsCount = float64(metrics.ParticipantsCount) / float64(metrics.RunningConferencesCount)
	} else {
		metrics.AverageParticipantsCount = 0
	}
	return metrics, nil
}

type TrendEntry struct {
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
}

func (r *ConferenceOrganizerRepo) GetParticipantsTrend(organizerId int) ([]TrendEntry, error) {
	var counts []TrendEntry

	var reportStartDate *time.Time
	query := `
	  SELECT MIN(c.start_date)
	  FROM ` + common.GetTableName(common.ConferenceTable) + ` c
	  JOIN ` + common.GetTableName(common.ConferenceOrganizerTable) + ` o ON c.id = o.conference_id
	  WHERE o.user_id = $1 AND c.start_date <= NOW() AND c.end_date >= NOW()
  `
	err := r.Db.Get(&reportStartDate, query, organizerId)
	if err != nil {
		return counts, errors.New("Could not get report start date")
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
	  LEFT JOIN ` + common.GetTableName(common.ConferenceParticipantTable) + ` cp
	  	ON cp.joined_at >= i.start_time
	  	AND cp.joined_at < i.start_time + $2::interval
	  LEFT JOIN ` + common.GetTableName(common.ConferenceOrganizerTable) + ` o
	  	ON cp.conference_id = o.conference_id
	  WHERE o.user_id = $3 OR cp.conference_id IS NULL
	  GROUP BY i.start_time
	  ORDER BY i.start_time ASC
  `

	err = r.Db.Select(&counts, query, reportStartDate, interval.String(), organizerId)
	if err != nil {
		return counts, errors.New("Could not get participants trend")
	}

	return counts, nil
}

func (r *ConferenceOrganizerRepo) GetConferenceOrganizerId(conferenceId int) (int, error) {
	var organizerId int
	query := "SELECT user_id FROM " + common.GetTableName(common.ConferenceOrganizerTable) + " WHERE conference_id = $1"
	err := r.Db.Get(
		&organizerId,
		query,
		conferenceId,
	)
	if err != nil {
		return 0, errors.New("Could not get organizer")
	}
	return organizerId, nil
}
