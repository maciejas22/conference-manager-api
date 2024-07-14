package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/directives"
	"github.com/maciejas22/conference-manager/api/graph"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/resolvers"
)

const defaultPort = "8080"

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "ResponseWriter", w))
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	db, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	r.Use(Middleware())
	r.Use(auth.Middleware())

	resolver := resolvers.NewResolver(ctx, db)
	c := graph.Config{Resolvers: resolver}
	c.Directives.Authenticated = directives.Authenticated
	c.Directives.HasRole = directives.HasRole
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	srv.AroundOperations(
		func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			rc := graphql.GetOperationContext(ctx)
			operationName := rc.Operation.Name
			startTime := time.Now()

			response := next(ctx)

			executionTime := time.Since(startTime)

			log.Printf(
				"Operation: %s took %v",
				operationName,
				executionTime,
			)
			return response
		},
	)

	r.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
