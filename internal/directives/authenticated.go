package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/maciejas22/conference-manager/api/internal/auth"
)

func Authenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthenticated")
	}

	return next(ctx)
}
