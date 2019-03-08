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
	yaml "gopkg.in/yaml.v2"
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

	reviewdir := os.Getenv("reviewdir")
	if reviewdir == "" {
		log.Fatal("$reviewdir must be set")
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
				// allows us to do review on monday/tuesday instead of just sunday
				// XXX: this year, week has to match up with logic in DumpSummaryForWeek
				year, week := time.Now().Add(-time.Hour * 72).ISOWeek()
				out, err := os.Create(fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week))
				if err != nil {
					log.Fatal(err)
				}

				err = workflow.DumpSummaryForWeek(user, appkey, authtoken, year, week, out)
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
				var inReviews []workflow.WeeklyReviewData

				for week := 1; week <= 53; week++ {
					weekSummary := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, year, week)
					if _, err1 := os.Stat(weekSummary); os.IsNotExist(err1) {
						// log.Printf("%+v doesn't exist, skipping", weekSummary)
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

					// get month
					var weekly workflow.WeeklySummary
					err3 := yaml.Unmarshal(buf, &weekly)
					if err3 != nil {
						log.Fatal(err3)
					}

					reviewFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, year, week)
					if _, err = os.Stat(reviewFname); os.IsNotExist(err) {
						log.Printf("%+v doesn't exist, skipping", reviewFname)
						continue
					}

					inReview, err2 := os.Open(reviewFname)
					if err2 != nil {
						log.Fatal(err2)
						continue
					}
					buf2, err2 := ioutil.ReadAll(inReview)
					if err2 != nil {
						log.Fatal(err2)
					}

					inReviews = append(inReviews, workflow.WeeklyReviewData{
						Week:    week,
						Month:   weekly.Month,
						Year:    weekly.Year,
						Content: buf2,
					})

				}

				for month := 1; month <= 12; month++ {
					monthlySummary := fmt.Sprintf("%s/monthly-%d-%02d.yaml", summarydir, year, month)
					out, err := os.Create(monthlySummary)
					if err != nil {
						log.Fatal(err)
					}
					err = workflow.GenerateSummaryForMonth(user, appkey, authtoken, year, month, inSummaries, inReviews, out)
					if err != nil {
						log.Printf("removing %+v, %+v", monthlySummary, err)
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
