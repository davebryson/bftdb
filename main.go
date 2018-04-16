package main

import (
	"log"
	"os"

	"github.com/davebryson/bftdb/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "bftdb cli"
	app.Version = "1.0"
	app.Description = "Tendermint + SQLite3"
	app.Author = "Dave Bryson"
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start the node",
			Action: func(c *cli.Context) error {
				cmd.RunNode()
				return nil
			},
		},
		{
			Name:  "console",
			Usage: "Run the interactive console",
			Action: func(c *cli.Context) error {
				cmd.RunConsole()
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
