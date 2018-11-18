package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/clarsen/go-trello-workflow/workflow"
)

var (
	appkey    string
	authtoken string
	user      string
	userEmail string
)

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
func HandleRequest(ctx context.Context, event myEvent) (string, error) {
	if event.Action == "minutely" {
		workflow.MinutelyMaintenance(user, appkey, authtoken)
	} else if event.Action == "daily" {
		workflow.DailyMaintenance(user, appkey, authtoken)
	}
	log.Printf("Handled event %s", event.Action)
	return fmt.Sprintf("Handled event %s!", event.Action), nil
}

func main() {
	lambda.Start(HandleRequest)
}
