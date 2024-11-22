package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	middlewares "github.com/maciejas22/conference-manager-api/cm-gateway/internal/middleware"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Authenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	si := middlewares.GetSessionInfo(ctx)
	if si.UserId == 0 {
		return nil, gqlerror.Errorf("Unauthenticated")
	}

	return next(ctx)
}
