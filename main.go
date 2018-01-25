package main

import (
	"log"
	"os"
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/robfig/cron"
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
	c := cron.New()
	// every night at 9:30PM
	c.AddFunc("0 30 21 * * *", func() { workflow.DailyMaintenance(user, appkey, authtoken) })
	// every minute
	c.AddFunc("0 * * * * *", func() { workflow.MinutelyMaintenance(user, appkey, authtoken) })
	c.Start()
	for {
		log.Println("wait 10 minutes...")
		time.Sleep(10 * time.Minute)
	}
}
