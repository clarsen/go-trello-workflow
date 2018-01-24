package main

import (
	"log"
	"os"

	"github.com/clarsen/go-trello-workflow/workflow"
)

func main() {
	appkey := os.Getenv("appkey")
	if appkey == "" {
		log.Fatal("$appkey must be set")
	}
	authtoken := os.Getenv("authtoken")
	if authtoken == "" {
		log.Fatal("$authtoken must be set")
	}
	user := os.Getenv("user")
	if authtoken == "" {
		log.Fatal("$authtoken must be set")
	}
	c, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
		return
	}
	c.Test()
}
