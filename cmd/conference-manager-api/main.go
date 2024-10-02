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

	"github.com/maciejas22/conference-manager/api/internal/config"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/directives"
	"github.com/maciejas22/conference-manager/api/internal/graph"
	"github.com/maciejas22/conference-manager/api/internal/middlewares"
	"github.com/maciejas22/conference-manager/api/internal/resolvers"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func initializeLogger(config *config.Config) *slog.Logger {
	var level slog.Level
	if config.GoEnv == "dev" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	return logger
}

func main() {
	ctx := context.Background()
	config.LoadConfig()

	logger := initializeLogger(config.AppConfig)

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedOrigins:   config.AppConfig.CorsAllowedOrigins,
		AllowCredentials: true,
		Debug:            config.AppConfig.GoEnv == "dev",
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

	resolver := resolvers.NewResolver(ctx, db, s3)
	c := graph.Config{Resolvers: resolver}
	c.Directives.Authenticated = directives.Authenticated
	c.Directives.HasRole = directives.HasRole
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	middlewares.LoadMiddlewares(r, srv, db, logger)

	port := config.AppConfig.Port
	r.Handle("/graphql", srv)
	if config.AppConfig.GoEnv == "dev" {
		r.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
		log.Printf("connect to http://localhost:%d/playground for GraphQL playground", port)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), r))
}
