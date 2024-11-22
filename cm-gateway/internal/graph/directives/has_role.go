package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/graph/model"
	middlewares "github.com/maciejas22/conference-manager-api/cm-gateway/internal/middleware"
	authPb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
	si := middlewares.GetSessionInfo(ctx)
	var userRole model.Role
	if si.Role == authPb.Role_ROLE_ORGANIZER {
		userRole = model.RoleOrganizer
	} else {
		userRole = model.RoleParticipant
	}
	if userRole != role {
		return nil, gqlerror.Errorf("You do not have the required role to perform this operation. Required role: %s", role)
	}

	return next(ctx)
}
