package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

func LoggingMiddleware(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Debug("http", "method", r.Method, "uri", r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

func GraphqlLogginMiddleware(srv *handler.Server, logger *slog.Logger) {
	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		req := graphql.GetOperationContext(ctx)
		logger.Debug("gql", "operation", req.Operation.Operation, "name", req.Operation.Name)
		return next(ctx)
	})
}
