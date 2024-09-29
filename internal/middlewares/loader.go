package middlewares

import (
	"log/slog"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi"
	"github.com/maciejas22/conference-manager/api/db"
)

func LoadMiddlewares(r *chi.Mux, srv *handler.Server, db *db.DB, l *slog.Logger) {
	r.Use(AuthMiddleware(db))
	r.Use(LoggingMiddleware(l))
	GraphqlLogginMiddleware(srv, l)
}
