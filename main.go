package main

import (
	"log"
	"os"
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/robfig/cron"
)

func updateBoards(user, appkey, authtoken string) {
	log.Println("Running updateBoards")
	wf, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
		return
	}
	wf.Test()
	log.Println("Finished running updateBoards")
}

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
	c := cron.New()
	c.AddFunc("@every 1m", func() { updateBoards(user, appkey, authtoken) })
	c.Start()
	for {
		log.Println("wait...")
		time.Sleep(60 * time.Second)
	}
}
