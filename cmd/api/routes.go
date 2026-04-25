package main

import (
	"net/http"

	"backend/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Post("/auth/google", app.googleAuth)

	mux.Group(func(r chi.Router) {
		r.Use(app.requireAuth)

		gqlSrv := handler.New(graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{DB: app.DB},
		}))
		gqlSrv.AddTransport(transport.POST{})
		gqlSrv.AddTransport(transport.GET{})
		gqlSrv.AddTransport(transport.Options{})

		r.Handle("/query", gqlSrv)
	})

	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/query"))

	return mux
}
