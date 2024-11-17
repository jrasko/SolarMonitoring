package main

import (
	"log"
	"net/http"
	"os"
	"pv-service/database"
	"pv-service/graph"
	"pv-service/graph/generated"
	"pv-service/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const port = "8080"

func main() {
	cfg := loadConfig()
	dbConnection, err := database.GetDBConnection(cfg.DBHost, cfg.DBUser, cfg.DBPassword)
	if err != nil {
		log.Fatalf("Error on startup has occured: %v", err)
	}

	processor := processing.New(dbConnection)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Processor: processor,
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/solar/query"))
	http.Handle("/solar/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
}

func loadConfig() Config {
	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
