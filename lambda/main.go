package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/clarsen/go-trello-workflow/workflow"
)

var (
	appkey    string
	authtoken string
	user      string
	userEmail string
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func init() {
	appkey = os.Getenv("appkey")
	if appkey == "" {
		log.Fatal("$appkey must be set")
	}
	authtoken = os.Getenv("authtoken")
	if authtoken == "" {
		log.Fatal("$authtoken must be set")
	}
	user = os.Getenv("user")
	if authtoken == "" {
		log.Fatal("$authtoken must be set")
	}
	userEmail = os.Getenv("USER_EMAIL")
	if userEmail == "" {
		log.Fatal("$USER_EMAIL must be set")
	}

}

type myEvent struct {
	Action string `json:"action"`
}

// HandleRequest accepts the lambda event
func HandleRequest(ctx context.Context, event myEvent) (Response, error) {
	var err error
	if event.Action == "minutely" {
		err = workflow.MinutelyMaintenance(user, appkey, authtoken)
	} else if event.Action == "daily" {
		err = workflow.DailyMaintenance(user, appkey, authtoken)
	} else if event.Action == "morning-reminder" {
		err = workflow.MorningRemind(user, appkey, authtoken, "", userEmail)
	} else if event.Action == "test" {
		resp1, err := http.Get("http://example.com/")
		if err != nil {
			return Response{StatusCode: 404}, err
		}
		_, err = ioutil.ReadAll(resp1.Body)
		resp1.Body.Close()
		if err != nil {
			return Response{StatusCode: 404}, err
		}

	}
	if err != nil {
		log.Printf("Handled event %s with error %+v", event.Action, err)
		resp := Response{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            fmt.Sprintf("Handled event %s with error %+v!", event.Action, err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}
		return resp, nil
	}
	log.Printf("Handled event %s", event.Action)
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            fmt.Sprintf("Handled event %s!", event.Action),
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
