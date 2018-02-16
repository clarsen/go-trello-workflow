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

	summarydir := os.Getenv("summarydir")
	if summarydir == "" {
		log.Fatal("$summarydir must be set")
	}

	reviewdir := os.Getenv("reviewdir")
	if reviewdir == "" {
		log.Fatal("$reviewdir must be set")
	}

	year, week := time.Now().ISOWeek()
	inSummary, err := os.Open(fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week))
	if err != nil {
		log.Fatal(err)
	}

	inReview, err := os.Open(fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, year, week))
	if err != nil {
		log.Fatal(err)
	}

	err = workflow.VisualizeWeeklyRetrospective(inSummary, inReview)
	if err != nil {
		log.Fatal(err)
	}

}
