package converters

import (
	"log"
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertTrendEntryRepoToSchema(c *repositories.TrendEntry) *models.ChartTrend {
	log.Println("Converting TrendEntryRepo to ChartTrend")
	log.Println("TrendEntryRepo: ", c)
	date, err := time.Parse(time.RFC3339, c.Date)
	log.Println("Date: ", date)
	if err != nil {
		return &models.ChartTrend{}
	}

	return &models.ChartTrend{
		Date:  date,
		Count: c.Value,
	}
}
