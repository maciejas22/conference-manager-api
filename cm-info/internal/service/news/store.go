package service

import (
	repoFilter "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/model"
	newsRepository "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/news"
	filter "github.com/maciejas22/conference-manager-api/cm-info/internal/service/model"
)

type NewsService struct {
	newsRepo newsRepository.NewsRepoInterface
}

func NewNewsService(newsRepo newsRepository.NewsRepoInterface) NewsServiceInterface {
	return &NewsService{
		newsRepo: newsRepo,
	}
}

func (s *NewsService) GetNews(page filter.Page) ([]News, filter.PaginationMeta, error) {
	pageFilter := repoFilter.Page{
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
	}
	news, meta, err := s.newsRepo.GetNews(pageFilter)
	if err != nil {
		return nil, filter.PaginationMeta{}, err
	}

	newsSlice := make([]News, len(news))
	for i, n := range news {
		newsSlice[i] = News{
			Id:        n.Id,
			Title:     n.Title,
			Content:   n.Content,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		}
	}
	newsMeta := filter.PaginationMeta{
		PageNumber: meta.PageNumber,
		PageSize:   meta.PageSize,
		TotalItems: meta.TotalItems,
		TotalPages: meta.TotalPages,
	}
	return newsSlice, newsMeta, nil
}
