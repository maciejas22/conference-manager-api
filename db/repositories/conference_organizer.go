package repositories

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type ConferenceOrganizer struct {
	UserId       string `json:"user_id" db:"user_id"`
	ConferenceId string `json:"conference_id" db:"conference_id"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

func (c *ConferenceOrganizer) TableName() string {
	return "public.conference_organizers"
}

func GetConferenceOrganizers(tx *sqlx.Tx, conferenceId string) ([]ConferenceOrganizer, error) {
	var organizers []ConferenceOrganizer
	o := &ConferenceOrganizer{}
	query := "SELECT user_id, conference_id FROM " + o.TableName() + " WHERE conference_id = $1"
	err := tx.Select(
		&organizers,
		query,
		conferenceId,
	)
	if err != nil {
		return nil, err
	}
	return organizers, nil
}

func IsConferenceOrganizer(tx *sqlx.Tx, conferenceId string, userId string) (bool, error) {
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

func AddConferenceOrganizer(tx *sqlx.Tx, conferenceId string, userId string) (Conference, error) {
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

func GetOrganizerLevelMetrics(tx *sqlx.Tx, organizerId string) (ConferenceOrganizerMetrics, error) {
	var metrics ConferenceOrganizerMetrics
	o := &ConferenceOrganizer{}
	c := &Conference{}
	cp := &ConferenceParticipant{}
	query := `SELECT 
					(SELECT COUNT(*) FROM ` + c.TableName() + ` WHERE id IN (SELECT conference_id FROM ` + o.TableName() + ` WHERE user_id = $1) AND NOW() BETWEEN start_date AND end_date) AS running_conferences_count,
					(SELECT COUNT(*) FROM ` + cp.TableName() + ` WHERE conference_id IN (SELECT id FROM ` + c.TableName() + ` WHERE id IN (SELECT conference_id FROM ` + o.TableName() + ` WHERE user_id = $1))) AS participants_count,
					(SELECT COUNT(*) FROM ` + o.TableName() + ` WHERE user_id = $1) AS total_organized_conferences`
	err := tx.Get(
		&metrics,
		query,
		organizerId,
	)
	if err != nil {
		log.Println("Error getting organizer metrics: ", err)
		return ConferenceOrganizerMetrics{}, err
	}
	metrics.AverageParticipantsCount = float64(metrics.ParticipantsCount) / float64(metrics.RunningConferencesCount)
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

func GetParticipantsTrend(tx *sqlx.Tx, organizerId string) (*ParticipantsTrend, error) {
	var counts []TrendEntry
	granularity := "Daily"

	var latestConference Conference
	query := `
		SELECT c.*
		FROM public.conferences c
		JOIN public.conference_organizers o ON c.id = o.conference_id
		WHERE o.user_id = $1 AND c.start_date <= NOW() AND c.end_date >= NOW()
		ORDER BY c.start_date DESC
		LIMIT 1
	`
	err := tx.Get(&latestConference, query, organizerId)
	if err != nil {
		log.Println("Error getting latest conference: ", err)
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
			WHERE conference_id = $1 AND joined_at >= $2 AND joined_at < $3
		`
		nextTime := t.Add(interval)
		err := tx.Get(&count, query, latestConference.Id, t, nextTime)
		if err != nil {
			log.Println("Error getting participants count: ", err)
			return nil, err
		}

		counts = append(counts, TrendEntry{
			Date:  t.Format(time.RFC3339),
			Value: count,
		})
	}
	log.Println("counts", counts)

	return &ParticipantsTrend{
		Trend:       counts,
		Granularity: granularity,
	}, nil
}
