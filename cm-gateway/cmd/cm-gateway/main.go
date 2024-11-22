package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/maciejas22/conference-manager-api/cm-gateway/config"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/client"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/graph"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/graph/directives"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/graph/resolvers"
	handlers "github.com/maciejas22/conference-manager-api/cm-gateway/internal/handler"
	middlewares "github.com/maciejas22/conference-manager-api/cm-gateway/internal/middleware"
)

func initLogger(config *config.Config) *slog.Logger {
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
	config.Init()

	logger := initLogger(config.AppConfig)

	client.SetupStripe(logger)
	authService := client.InitAuthClient(ctx)
	conferenceService := client.InitConferencesClient(ctx)

	resolver := resolvers.NewResolver(ctx, logger, authService, conferenceService)
	c := graph.Config{Resolvers: resolver}
	c.Directives.Authenticated = directives.Authenticated
	c.Directives.HasRole = directives.HasRole
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	port := strconv.Itoa(config.AppConfig.Port)
	if config.AppConfig.GoEnv == config.EnvDev {
		logger.Info("gql", "serving playground url", fmt.Sprintf("http://localhost:%s/playground", port))
		http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	}

	middlewares.LoadGqlMiddlewares(srv, logger)
	m := middlewares.LoadHttpMiddlewares(logger, authService)
	logger.Info("gql", "serving graphql url", fmt.Sprintf("http://localhost:%s/graphql", port))
	http.Handle("/graphql", m(srv))

	logger.Info("stripe", "url", fmt.Sprintf("http://localhost:%s/v1/webhooks/stripe", port))
	http.HandleFunc("/v1/webhooks/stripe", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePaymentIntentConfirmation(ctx, logger, conferenceService)(w, r)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
