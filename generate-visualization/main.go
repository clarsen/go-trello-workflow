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

	reviewvisdir := os.Getenv("reviewvisdir")
	if reviewvisdir == "" {
		log.Fatal("$reviewvisdir must be set")
	}

	app := cli.NewApp()
	app.Name = "trello-dump-summary"
	app.Usage = "Dump data from Trello board"
	app.Commands = []cli.Command{
		{
			Name:    "weekly-review-template",
			Aliases: []string{"tw"},
			Usage:   "Generate weekly review \"visualization\" template",
			Action: func(*cli.Context) {

				year, week := time.Now().AddDate(0, 0, -3).ISOWeek()

				templateFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, year, week)
				if _, err = os.Stat(templateFname); err == nil {
					log.Fatalf("%s exists already", templateFname)
				}

				outReview, err2 := os.Create(templateFname)
				if err2 != nil {
					log.Fatal(err2)
				}

				err2 = workflow.CreateEmptyWeeklyRetrospective(outReview)
				if err2 != nil {
					log.Fatal(err2)
				}
			},
		},
		{
			Name:    "weekly-review",
			Aliases: []string{"w"},
			Usage:   "Generate weekly review \"visualization\"",
			Action: func(*cli.Context) {
				year := time.Now().AddDate(0, 0, -3).Year()
				for week := 1; week <= 53; week++ {
					summaryFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week)
					if _, err = os.Stat(summaryFname); os.IsNotExist(err) {
						continue
					}

					inSummary, err2 := os.Open(summaryFname)
					if err2 != nil {
						log.Fatal(err2)
					}
					reviewFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, year, week)
					if _, err = os.Stat(reviewFname); os.IsNotExist(err) {
						continue
					}

					inReview, err2 := os.Open(reviewFname)
					if err2 != nil {
						log.Fatal(err2)
						continue
					}

					out, err2 := os.Create(fmt.Sprintf("%s/weekly-%d-%02d.md", reviewvisdir, year, week))
					if err2 != nil {
						log.Fatal(err2)
					}

					err2 = workflow.VisualizeWeeklyRetrospective(inSummary, inReview, out)
					if err2 != nil {
						log.Fatal(err2)
					}
				}
			},
		},
		{
			Name:    "monthly-review-template",
			Aliases: []string{"tm"},
			Usage:   "Generate monthly review \"visualization\" template",
			Action: func(*cli.Context) {

				year := time.Now().AddDate(0, 0, -3).Year()
				month := int(time.Now().AddDate(0, 0, -5).Month())

				inSummary, err2 := os.Open(fmt.Sprintf("%s/monthly-%d-%02d.yaml", summarydir, year, month))
				if err != nil {
					log.Fatal(err2)
				}

				templateFname := fmt.Sprintf("%s/monthly-%d-%02d.yaml", reviewdir, year, month)
				if _, err = os.Stat(templateFname); err == nil {
					log.Fatalf("%s exists already", templateFname)
				}

				outReview, err2 := os.Create(templateFname)
				if err2 != nil {
					log.Fatal(err2)
				}

				err2 = workflow.CreateEmptyMonthlyRetrospective(inSummary, outReview)
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

				for month := 1; month <= 12; month++ {
					inSummary, err := os.Open(fmt.Sprintf("%s/monthly-%d-%02d.yaml", summarydir, year, month))
					if err != nil {
						log.Println(err)
						continue
					}

					inReview, err := os.Open(fmt.Sprintf("%s/monthly-%d-%02d.yaml", reviewdir, year, month))
					if err != nil {
						log.Println(err)
						continue
					}
					out, err := os.Create(fmt.Sprintf("%s/monthly-%d-%02d.md", reviewvisdir, year, month))
					if err != nil {
						log.Fatal(err)
					}

					err = workflow.VisualizeMonthlyRetrospective(inSummary, inReview, out)
					if err != nil {
						log.Fatal(err)
					}

				}

			},
		},
	}
	app.Run(os.Args)

}
