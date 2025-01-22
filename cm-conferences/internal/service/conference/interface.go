package service

import (
	"context"

	common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/common"
)

type ConferenceServiceInterface interface {
	CreateConference(ctx context.Context, userId int, createConferenceInput CreateConferenceInput) (*int, error)
	ModifyConference(ctx context.Context, input ModifyConferenceInput) (*int, error)
	GetConferencesPage(ctx context.Context, userId int, p *common.Page, sort *common.Sort, f *ConferencesFilters) ([]int, common.PaginationMeta, error)
	GetConference(ctx context.Context, id int) (Conference, error)
	GetConferencesMetrics(ctx context.Context) (ConferencesMetrics, error)
	GetAgenda(ctx context.Context, conferenceId int) ([]AgendaItem, error)
	GetConferencesByIds(ctx context.Context, ids []int) ([]Conference, error)
}
