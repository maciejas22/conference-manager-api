package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role models.Role) (interface{}, error) {
	si := auth.GetSessionInfo(ctx)
	if models.Role(si.Role) != role {
		return nil, gqlerror.Errorf("You do not have the required role to perform this operation. Required role: %s", role)
	}

	return next(ctx)
}
