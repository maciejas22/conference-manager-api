package converters

import (
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertUserRepoToSchema(u *repositories.User) *models.User {
	return &models.User{
		ID:       u.Id,
		Name:     u.Name,
		Surname:  u.Surname,
		Username: u.Username,
		Email:    u.Email,
		Role:     models.Role(u.Role),
	}
}
