package resolvers

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
)

type Resolver struct {
	ctx                       context.Context
	agendaRepo                repositories.AgendaRepository
	conferenceOrganizerRepo   repositories.ConferenceOrganizerRepository
	conferenceParticipantRepo repositories.ConferenceParticipantRepository
	conferenceRepo            repositories.ConferenceRepository
	newsRepo                  repositories.NewsRepository
	sectionRepo               repositories.SectionRepository
	subsectionRepo            repositories.SubsectionRepository
	termsOfServiceRepo        repositories.TermsOfServiceRepository
	userRepo                  repositories.UserRepository
}

func NewResolver(ctx context.Context, db *db.DB) *Resolver {
	return &Resolver{
		ctx:                       ctx,
		agendaRepo:                repositories.NewAgendaRepository(ctx, db),
		conferenceOrganizerRepo:   repositories.NewConferenceOrganizerRepository(ctx, db),
		conferenceParticipantRepo: repositories.NewConferenceParticipantRepository(ctx, db),
		conferenceRepo:            repositories.NewConferenceRepository(ctx, db),
		newsRepo:                  repositories.NewNewsRepository(ctx, db),
		sectionRepo:               repositories.NewSectionRepository(ctx, db),
		subsectionRepo:            repositories.NewSubsectionRepository(ctx, db),
		termsOfServiceRepo:        repositories.NewTermsOfServiceRepository(ctx, db),
		userRepo:                  repositories.NewUserRepo(ctx, db),
	}
}
