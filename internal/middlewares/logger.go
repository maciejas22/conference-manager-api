package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

func LoggingMiddleware(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logMessage := fmt.Sprintf("Method: %s, URI: %s, User Agent: %s, Remote Addr: %s, Time: %s", r.Method, r.RequestURI, r.UserAgent(), r.RemoteAddr, start.Format(time.RFC1123))
			l.Debug(logMessage)
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
