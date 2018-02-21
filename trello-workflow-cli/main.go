package main

import (
	"log"
	"os"

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
	if user == "" {
		log.Fatal("$user must be set")
	}

	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		log.Fatal("$sendgridKey must be set")
	}

	userEmail := os.Getenv("USER_EMAIL")
	if userEmail == "" {
		log.Fatal("$USER_EMAIL must be set")

	}

	app := cli.NewApp()
	app.Name = "trello-workflow"
	app.Usage = "Update Trello board"
	app.Commands = []cli.Command{
		{
			Name:    "remind",
			Aliases: []string{"r"},
			Usage:   "Update the trello board on daily basis",
			Action: func(*cli.Context) {
				err := workflow.MorningRemind(user, appkey, authtoken, sendgridKey, userEmail)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "today",
			Aliases: []string{"t"},
			Usage:   "Update the trello board on daily basis",
			Action:  func(*cli.Context) { workflow.DailyMaintenance(user, appkey, authtoken) },
		},
		{
			Name:    "minutely",
			Aliases: []string{"m"},
			Usage:   "Update the trello board on minutely basis",
			Action:  func(*cli.Context) { workflow.MinutelyMaintenance(user, appkey, authtoken) },
		},
		{
			Name:    "weekly",
			Aliases: []string{"w"},
			Usage:   "End of week report",
			Action: func(*cli.Context) {
				err := workflow.Weekly(user, appkey, authtoken)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "weekly cleanup",
			Aliases: []string{"wc"},
			Usage:   "End of week cleanup",
			Action: func(*cli.Context) {
				err := workflow.WeeklyCleanup(user, appkey, authtoken)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)

}
