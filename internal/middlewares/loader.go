package middlewares

import "github.com/go-chi/chi"

func LoadMiddlewares(r *chi.Mux) {
	r.Use(AuthMiddleware())
	r.Use(LoggingMiddleware)
}
