package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/clarsen/go-trello-workflow/server/go/handle_graphql"
)

var muxAdapter *gorillamux.GorillaMuxAdapter

func addCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Required for CORS support to work
		w.Header().Set("Access-Control-Allow-Origin", "https://enchilada-serverless-next-auth0.app.caselarsen.com")
		// Required for cookies, authorization headers with HTTPS
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func lambdahandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("got request %+v\n", request)
	res, err := muxAdapter.Proxy(request)
	log.Printf("return response %+v\n", res)
	if err != nil {
		log.Println(err)
	}
	return res, err

}

func init() {
	// start the mux router
	r := mux.NewRouter()
	r.Use(addCors)
	r.Use(handle_graphql.GetAuthID)
	r.HandleFunc("/api/gql", handler.GraphQL(handle_graphql.NewExecutableSchema(handle_graphql.Config{Resolvers: &handle_graphql.Resolver{}})))
	// api.Routes(r) // routes just accepts a mux router to bind routes to and that's where our gqlgen handlers live
	// initialize the mux adapter so that we can use mux with lambda
	muxAdapter = gorillamux.New(r)
}

func main() {
	lambda.Start(lambdahandler)
}
