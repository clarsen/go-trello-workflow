package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
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

	app := cli.NewApp()
	app.Name = "trello-dump-summary"
	app.Usage = "Dump data from Trello board"
	app.Commands = []cli.Command{
		{
			Name:    "weekly-review",
			Aliases: []string{"w"},
			Usage:   "Generate weekly review \"visualization\"",
			Action: func(*cli.Context) {
				year, week := time.Now().AddDate(0, 0, -3).ISOWeek()
				inSummary, err2 := os.Open(fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week))
				if err != nil {
					log.Fatal(err2)
				}

				inReview, err2 := os.Open(fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, year, week))
				if err != nil {
					log.Fatal(err2)
				}

				err2 = workflow.VisualizeWeeklyRetrospective(inSummary, inReview)
				if err2 != nil {
					log.Fatal(err2)
				}
			},
		},
		{
			Name:    "monthly-review",
			Aliases: []string{"m"},
			Usage:   "Generate monthly review \"visualization\"",
			Action: func(*cli.Context) {
				year := time.Now().Year()
				month := int(time.Now().AddDate(0, 0, -5).Month())

				inSummary, err := os.Open(fmt.Sprintf("%s/monthly-%d-%02d.yaml", summarydir, year, month))
				if err != nil {
					log.Fatal(err)
				}

				inReview, err := os.Open(fmt.Sprintf("%s/monthly-%d-%02d.yaml", reviewdir, year, month))
				if err != nil {
					log.Fatal(err)
				}

				err = workflow.VisualizeMonthlyRetrospective(inSummary, inReview)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)

}
