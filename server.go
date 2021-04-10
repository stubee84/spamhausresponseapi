package main

import (
	"SpamhausResponseApi/db"
	"SpamhausResponseApi/graph"
	"SpamhausResponseApi/graph/generated"
	"SpamhausResponseApi/helpers"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Missing the PORT environment variable")
	}

	db.Connect()

	router := chi.NewRouter()
	router.Route("/graphql", func(r chi.Router) {
		r.Use(helpers.BasicAuth())
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
		r.Handle("/", srv)
	})

	log.Printf("connect to http://localhost:%s/graphql for GraphQL queries", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
