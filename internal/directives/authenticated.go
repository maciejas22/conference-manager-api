package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Authenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	si := auth.GetSessionInfo(ctx)
	if si.UserID == 0 {
		return nil, gqlerror.Errorf("Unauthenticated")
	}

	return next(ctx)
}
