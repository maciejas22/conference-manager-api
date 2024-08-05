package converter

import (
	filters "github.com/maciejas22/conference-manager/api/db/repositories/shared"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertOrderSchemaToRepo(o models.Order) filters.Order {
	if o == models.OrderAsc {
		return filters.ASC
	}
	return filters.DESC
}

func ConvertSortSchemaToRepo(s *models.Sort) *filters.Sort {
	if s == nil {
		return nil
	}

	return &filters.Sort{
		Column: s.Column,
		Order:  ConvertOrderSchemaToRepo(s.Order),
	}
}
