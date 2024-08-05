package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/internal/config"
	"github.com/maciejas22/conference-manager/api/internal/directives"
	"github.com/maciejas22/conference-manager/api/internal/graph"
	"github.com/maciejas22/conference-manager/api/internal/middlewares"
	"github.com/maciejas22/conference-manager/api/internal/resolvers"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx := context.Background()
	config.LoadConfig()

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	db, err := db.Connect(ctx, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
	}
	defer db.Close()

	s3, err := s3.NewS3Client(logger)
	if err != nil {
		logger.Error("failed to connect to s3", "error", err)
	}

	middlewares.LoadMiddlewares(r)

	resolver := resolvers.NewResolver(ctx, db, s3)
	c := graph.Config{Resolvers: resolver}
	c.Directives.Authenticated = directives.Authenticated
	c.Directives.HasRole = directives.HasRole
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", srv)

	port := config.AppConfig.Port
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), r))
}
