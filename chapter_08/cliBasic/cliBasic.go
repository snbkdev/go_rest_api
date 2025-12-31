package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "name",
			Value: "stranger",
			Usage: "your wonderful world",
		},
		cli.IntFlag{
			Name: "age",
			Value: 0,
			Usage: "your graceful age",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Printf("Hello %s (%d years), welcome to the command line world", c.String("name"), c.Int("age"))
		return nil
	}

	app.Run(os.Args)
}