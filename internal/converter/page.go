package converter

import (
	filters "github.com/maciejas22/conference-manager/api/db/repositories/shared"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertPageSchemaToRepo(p *models.Page) filters.Page {
	if p == nil {
		return filters.Page{
			PageNumber: 1,
			PageSize:   10,
		}
	}

	return filters.Page{
		PageNumber: p.Number,
		PageSize:   p.Size,
	}
}

func ConvertPageMetaToRepo(p *models.PageInfo) *filters.PaginationMeta {
	if p == nil {
		return nil
	}

	return &filters.PaginationMeta{
		PageNumber: p.Number,
		PageSize:   p.Size,
		TotalItems: p.TotalItems,
		TotalPages: p.TotalPages,
	}
}
