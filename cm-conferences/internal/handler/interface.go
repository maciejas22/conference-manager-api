package grpc

import (
	a "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/agenda"
	c "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/conference"
	cp "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/organizer"
	co "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/participant"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/conference"
)

type ConferenceServ struct {
	conferenceService  c.ConferenceServiceInterface
	organizerService   cp.OrganizerServiceInterface
	participantService co.ParticipantServiceInterface
	agendaService      a.AgendaServiceInterface
	pb.UnimplementedConferenceServiceServer
}
