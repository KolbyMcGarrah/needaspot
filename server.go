package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/KolbyMcGarrah/nas/graph"
	"github.com/KolbyMcGarrah/nas/graph/generated"
	"github.com/KolbyMcGarrah/nas/internal/auth"
	database "github.com/KolbyMcGarrah/nas/internal/pkg/db/postgres"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//initialize new router
	router := chi.NewRouter()

	//assign our middleware to the router
	router.Use(auth.Middleware())

	//initialize Postgres Db connection
	database.InitDB()
	//run any migrations that haven't been applied.
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
