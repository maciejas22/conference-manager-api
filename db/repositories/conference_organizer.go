package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
)

type ConferenceOrganizer struct {
	UserId       int    `json:"user_id" db:"user_id"`
	ConferenceId int    `json:"conference_id" db:"conference_id"`
	JoinedAt     string `json:"joined_at" db:"joined_at"`
}

func (c *ConferenceOrganizer) TableName() string {
	return "conference_organizers"
}

func GetConferenceOrganizers(qe *db.QueryExecutor, conferenceId int) ([]ConferenceOrganizer, error) {
	var organizers []ConferenceOrganizer
	o := &ConferenceOrganizer{}
	query := "SELECT user_id, conference_id FROM " + o.TableName() + " WHERE conference_id = ?"
	err := sqlx.Select(
		qe,
		&organizers,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return organizers, nil
}

func IsConferenceOrganizer(qe *db.QueryExecutor, conferenceId int, userId int) (bool, error) {
	var count int
	o := &ConferenceOrganizer{}
	query := "SELECT COUNT(*) FROM " + o.TableName() + " WHERE conference_id = ? AND user_id = ?"
	err := sqlx.Get(
		qe,
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

func AddConferenceOrganizer(qe *db.QueryExecutor, conferenceId int, userId int) (Conference, error) {
	o := &ConferenceOrganizer{
		UserId:       userId,
		ConferenceId: conferenceId,
	}
	query := "INSERT INTO " + o.TableName() + " (user_id, conference_id) VALUES (?, ?)"
	_, err := qe.Exec(
		query,
		o.UserId,
		o.ConferenceId,
	)
	if err != nil {
		return Conference{}, err
	}
	return Conference{}, nil
}

type ConferenceOrganizerMetrics struct {
	RunningConferencesCount   int     `db:"running_conferences_count"`
	ParticipantsCount         int     `db:"participants_count"`
	AverageParticipantsCount  float64 `db:"average_participants_count"`
	TotalOrganizedConferences int     `db:"total_organized_conferences"`
}

func GetOrganizerLevelMetrics(qe *db.QueryExecutor, organizerId int) (ConferenceOrganizerMetrics, error) {
	var metrics ConferenceOrganizerMetrics
	o := &ConferenceOrganizer{}
	c := &Conference{}
	cp := &ConferenceParticipant{}
	query := `SELECT 
					(SELECT COUNT(*) FROM ` + c.TableName() + ` WHERE id IN (SELECT conference_id FROM ` + o.TableName() + ` WHERE user_id = ?) AND datetime('now') BETWEEN start_date AND end_date) AS running_conferences_count,
					(SELECT COUNT(*) FROM ` + cp.TableName() + ` WHERE conference_id IN (SELECT id FROM ` + c.TableName() + ` WHERE id IN (SELECT conference_id FROM ` + o.TableName() + ` WHERE user_id = ?))) AS participants_count,
					(SELECT COUNT(*) FROM ` + o.TableName() + ` WHERE user_id = ?) AS total_organized_conferences`
	err := sqlx.Get(
		qe,
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
	Date  string
	Value int
}

type ParticipantsTrend struct {
	Trend       []TrendEntry
	Granularity string
}

func GetParticipantsTrend(qe *db.QueryExecutor, organizerId int) (*ParticipantsTrend, error) {
	var counts []TrendEntry
	granularity := "Daily"

	var latestConference Conference
	query := `
		SELECT c.*
		FROM ` + latestConference.TableName() + ` c
		JOIN ` + (new(ConferenceOrganizer)).TableName() + `o ON c.id = o.conference_id
		WHERE o.user_id = ? AND c.start_date <= datetime('now') AND c.end_date >= datetime('now')
		ORDER BY c.start_date DESC
		LIMIT 1
	`
	err := sqlx.Get(qe, &latestConference, query, organizerId)
	if err != nil {
		return &ParticipantsTrend{
			Trend:       counts,
			Granularity: granularity,
		}, nil
	}

	startTime, _ := time.Parse(time.RFC3339, latestConference.StartDate)
	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)

	var interval time.Duration
	if totalDuration.Hours() <= 10*24 {
		interval = 24 * time.Hour
		granularity = "Daily"
	} else if totalDuration.Hours() <= 10*24*7 {
		interval = 7 * 24 * time.Hour
		granularity = "Weekly"
	} else {
		interval = 30 * 24 * time.Hour
		granularity = "Monthly"
	}

	for t := startTime; t.Before(endTime); t = t.Add(interval) {
		var count int
		query = `
			SELECT COUNT(*)
			FROM public.conference_participants
			WHERE conference_id = ? AND joined_at >= ? AND joined_at < ? 
		`
		nextTime := t.Add(interval)
		err := sqlx.Get(qe, &count, query, latestConference.Id, t, nextTime)
		if err != nil {
			return nil, err
		}

		counts = append(counts, TrendEntry{
			Date:  t.Format(time.RFC3339),
			Value: count,
		})
	}

	return &ParticipantsTrend{
		Trend:       counts,
		Granularity: granularity,
	}, nil
}
