package main

import (
	"log"
	"net/http"
	"os"
	"pv-service/graph"
	"pv-service/graph/generated"
	"pv-service/processing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	processor, err := processing.GetProcessor()
	if err != nil {
		log.Printf("Error on startup has occured: %v", err)
		panic(err)
	}
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Processor: processor,
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
