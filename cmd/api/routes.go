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

	mux.Route("/users", func(r chi.Router) {
		r.Get("/", app.ListUsers)
		r.Post("/", app.CreateUser)
		r.Get("/{id}", app.GetUser)
		r.Put("/{id}", app.UpdateUser)
		r.Delete("/{id}", app.DeleteUser)
	})

	mux.Route("/worlds", func(r chi.Router) {
		r.Get("/", app.ListWorlds)
		r.Post("/", app.CreateWorld)
		r.Get("/{id}", app.GetWorld)
		r.Put("/{id}", app.UpdateWorld)
		r.Delete("/{id}", app.DeleteWorld)
	})

	mux.Route("/map-assets", func(r chi.Router) {
		r.Get("/", app.ListMapAssets)
		r.Post("/", app.CreateMapAsset)
		r.Get("/{id}", app.GetMapAsset)
		r.Put("/{id}", app.UpdateMapAsset)
		r.Delete("/{id}", app.DeleteMapAsset)
	})

	gqlSrv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{DB: app.DB},
	}))
	gqlSrv.AddTransport(transport.POST{})
	gqlSrv.AddTransport(transport.GET{})
	gqlSrv.AddTransport(transport.Options{})

	mux.Handle("/query", gqlSrv)
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/query"))

	return mux
}
