package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	authPb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
)

func chainHttpMiddlewares(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(finalHandler http.Handler) http.Handler {
		for _, middleware := range middlewares {
			finalHandler = middleware(finalHandler)
		}
		return finalHandler
	}
}

func LoadHttpMiddlewares(l *slog.Logger, authService authPb.AuthServiceClient) func(http.Handler) http.Handler {
	// logginMiddleware := LoggingMiddleware(l)
	authMiddleware := AuthMiddleware(authService)

	// chainedHandler := chainHttpMiddlewares(logginMiddleware, authMiddleware)
	chainedHandler := chainHttpMiddlewares(authMiddleware)
	return chainedHandler
}

func LoadGqlMiddlewares(srv *handler.Server, l *slog.Logger) {
	GraphqlLogginMiddleware(srv, l)
}
