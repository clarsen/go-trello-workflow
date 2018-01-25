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
	if authtoken == "" {
		log.Fatal("$authtoken must be set")
	}

	app := cli.NewApp()
	app.Name = "trello-workflow"
	app.Usage = "Update Trello board"
	app.Commands = []cli.Command{
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
	}
	app.Run(os.Args)

}
