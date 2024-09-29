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
	return "conference_organizers"
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

	var latestConference Conference
	query := `
		SELECT c.*
		FROM ` + latestConference.TableName() + ` c
		JOIN ` + (new(ConferenceOrganizer)).TableName() + ` o ON c.id = o.conference_id
		WHERE o.user_id = $1 AND c.start_date <= NOW() AND c.end_date >= NOW()
		ORDER BY c.start_date DESC
		LIMIT 1
	`
	err := tx.Get(&latestConference, query, organizerId)
	if err != nil {
		return counts, nil
	}

	startTime, _ := time.Parse(time.RFC3339, latestConference.StartDate)
	endTime := time.Now()
	totalDuration := endTime.Sub(startTime)

	var interval time.Duration
	if totalDuration.Hours() <= 24 {
		interval = 24 * time.Hour / 10
	} else {
		interval = totalDuration / 10
	}

	for t := startTime; t.Before(endTime); t = t.Add(interval) {
		var count int
		query = `
			SELECT COUNT(*)
			FROM public.conference_participants
			WHERE conference_id = $1 AND joined_at >= $2 AND joined_at < $3 
		`
		nextTime := t.Add(interval)
		err := tx.Get(&count, query, latestConference.Id, t, nextTime)
		if err != nil {
			return nil, err
		}

		counts = append(counts, TrendEntry{
			Date:  t.Format(time.RFC3339),
			Value: count,
		})
	}

	return counts, nil
}
