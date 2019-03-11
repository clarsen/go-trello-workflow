package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/handler"
	"github.com/clarsen/go-trello-workflow/server/go/handle_graphql"
	"github.com/gorilla/mux"
)

const defaultPort = "8080"

func addCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Required for CORS support to work
		w.Header().Set("Access-Control-Allow-Origin", "https://https://workflow.app.caselarsen.com")
		if strings.HasPrefix(r.Header.Get("Origin"), "http://localhost") {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Amz-Date,Authorization")
		// Required for cookies, authorization headers with HTTPS
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := mux.NewRouter()
	router.Use(addCors)
	router.Handle("/", handler.Playground("GraphQL playground", "/api/gql"))
	router.Handle("/api/gql",
		handler.GraphQL(handle_graphql.NewExecutableSchema(handle_graphql.Config{Resolvers: &handle_graphql.Resolver{}}))).Methods("POST", "OPTIONS")

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
