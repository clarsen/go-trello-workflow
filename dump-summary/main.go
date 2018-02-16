package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	summarydir := os.Getenv("summarydir")
	if summarydir == "" {
		log.Fatal("$summarydir must be set")
	}

	year, week := time.Now().ISOWeek()
	out, err := os.Create(fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week))
	if err != nil {
		log.Fatal(err)
	}

	err = workflow.DumpSummaryForWeek(user, appkey, authtoken, out)
	if err != nil {
		log.Fatal(err)
	}

}
