package service

import (
	a "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/agenda"
	c "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/conference"
	co "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/organizer"
	cp "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/participant"
)

type Services struct {
	ConferenceService  c.ConferenceServiceInterface
	OrganizerService   co.OrganizerServiceInterface
	ParticipantService cp.ParticipantServiceInterface
	AgendaService      a.AgendaServiceInterface
}

func InitServices(
	conferenceService c.ConferenceServiceInterface,
	organizerService co.OrganizerServiceInterface,
	participantService cp.ParticipantServiceInterface,
	agendaService a.AgendaServiceInterface,
) Services {
	return Services{
		ConferenceService:  conferenceService,
		OrganizerService:   organizerService,
		ParticipantService: participantService,
		AgendaService:      agendaService,
	}
}
