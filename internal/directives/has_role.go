package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role models.Role) (interface{}, error) {
	claims, ok := auth.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthenticated")
	}

	if models.Role(claims.UserMetadata.Role) != role {
		return nil, fmt.Errorf("required role: %s", role)
	}

	return next(ctx)
}
