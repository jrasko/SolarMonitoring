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

	// http.Handle("/", playground.Handler("GraphQL playground", "/solar/query"))
	http.Handle("/solar/query", basicAuth(srv, cfg.BasicUsername, cfg.BasicPassword))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string

	BasicUsername string
	BasicPassword string
}

func loadConfig() Config {
	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),

		BasicUsername: os.Getenv("BASIC_USERNAME"),
		BasicPassword: os.Getenv("BASIC_PASSWORD"),
	}
}

func basicAuth(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok && user == username && pass == password {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
