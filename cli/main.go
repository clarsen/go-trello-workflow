package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "trello-workflow"
	app.Usage = "Update Trello board"
	app.Commands = []cli.Command{
		{
			Name:    "today",
			Aliases: []string{"t"},
			Usage:   "Update the trello board on daily basis",
			Action:  doToday,
		},
	}
	app.Run(os.Args)

}
