package main

import (
	"fmt"
	"io/ioutil"
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

	app := cli.NewApp()
	app.Name = "trello-dump-summary"
	app.Usage = "Dump data from Trello board"
	app.Commands = []cli.Command{
		{
			Name:    "week",
			Aliases: []string{"w"},
			Usage:   "Summarize the week in progress or end of week",
			Action: func(*cli.Context) {
				year, week := time.Now().ISOWeek()
				out, err := os.Create(fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week))
				if err != nil {
					log.Fatal(err)
				}

				err = workflow.DumpSummaryForWeek(user, appkey, authtoken, out)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "historical",
			Aliases: []string{"h"},
			Usage:   "Summarize from history (not needed except for summarization schema update)",
			Action: func(*cli.Context) {
				year := 2018
				if year != time.Now().Year() {
					log.Fatal("doesn't make sense after 2018")
				}
				// for week D, dump week D
				for week := 1; week < 7; week++ {
					log.Println("week", week)
					out, err := os.Create(fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week))
					if err != nil {
						log.Fatal(err)
					}
					err = workflow.DumpSummaryForWeekFromHistory(user, appkey, authtoken, week, out)
					if err != nil {
						log.Fatal(err)
					}
					if week != 6 {
						log.Println("Waiting...")
						time.Sleep(time.Second * 30)
					}
				}
			},
		},
		{
			Name:    "month",
			Aliases: []string{"m"},
			Usage:   "Summarize the month in progress or end of month, and all previous months too",
			Action: func(*cli.Context) {
				year := time.Now().Year()
				var inSummaries [][]byte
				for week := 1; week <= 53; week++ {
					weekSummary := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week)
					if _, err := os.Stat(weekSummary); os.IsNotExist(err) {
						continue
					}
					inSummary, err2 := os.Open(weekSummary)
					if err2 != nil {
						log.Fatal(err2)
					}
					buf, err2 := ioutil.ReadAll(inSummary)
					if err2 != nil {
						log.Fatal(err2)
					}
					inSummaries = append(inSummaries, buf)
				}
				for month := 1; month <= 12; month++ {
					monthlySummary := fmt.Sprintf("%s/monthly-%d-%02d.yaml", summarydir, year, month)
					out, err := os.Create(monthlySummary)
					if err != nil {
						log.Fatal(err)
					}
					err = workflow.GenerateSummaryForMonth(user, appkey, authtoken, year, month, inSummaries, out)
					if err != nil {
						os.Remove(monthlySummary)
					}
				}
			},
		},
		{
			Name:    "year",
			Aliases: []string{"y"},
			Usage:   "Summarize the year in progress or end of year",
			Action: func(*cli.Context) {
				year := time.Now().Year()

				var inMonthlySummaries [][]byte
				for month := 1; month <= 12; month++ {
					monthlySummary := fmt.Sprintf("%s/monthly-%d-%02d.yaml", summarydir, year, month)
					if _, err := os.Stat(monthlySummary); os.IsNotExist(err) {
						continue
					}
					inSummary, err2 := os.Open(monthlySummary)
					if err2 != nil {
						log.Fatal(err2)
					}
					buf, err2 := ioutil.ReadAll(inSummary)
					if err2 != nil {
						log.Fatal(err2)
					}
					inMonthlySummaries = append(inMonthlySummaries, buf)
				}

				yearlySummary := fmt.Sprintf("%s/yearly-%d.yaml", summarydir, year)
				out, err := os.Create(yearlySummary)
				if err != nil {
					log.Fatal(err)
				}
				err = workflow.GenerateSummaryForYear(year, inMonthlySummaries, out)
				if err != nil {
					os.Remove(yearlySummary)
				}

			},
		},
	}
	app.Run(os.Args)

}
