package repository

import (
	"github.com/maciejas22/conference-manager-api/cm-conferences/internal/db"
	a "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/agenda"
	c "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference"
	co "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference_organizer"
	cp "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference_participant"
)

type Repos struct {
	ConferenceRepo  c.ConferenceRepoInterface
	OrganizerRepo   co.ConferenceOrganizerRepoInterface
	ParticipantRepo cp.ConferenceParticipantRepoInterface
	AgendaRepo      a.AgendaRepoInterface
}

func InitRepos(db *db.DB) Repos {
	return Repos{
		ConferenceRepo:  c.NewConferenceRepo(db.Conn),
		OrganizerRepo:   co.NewOrganizerRepo(db.Conn),
		ParticipantRepo: cp.NewParticipantRepo(db.Conn),
		AgendaRepo:      a.NewAgendaRepo(db.Conn),
	}
}
