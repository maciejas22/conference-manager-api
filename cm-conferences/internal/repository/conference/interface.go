package repository

import common "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/common"

type ConferenceRepoInterface interface {
	GetConferenceById(conferenceId int) (Conference, error)
	GetConferencesByIds(conferenceIds []int) ([]Conference, error)
	GetConferencesPage(userId int, p common.Page, s *common.Sort, f *ConferenceFilter) ([]int, common.PaginationMeta, error)
	CreateConference(conference Conference, organizerId int) (int, error)
	UpdateConference(conference Conference) (int, error)
	GetMetrics() (ConferencesMetrics, error)
}
