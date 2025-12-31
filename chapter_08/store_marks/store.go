package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "save",
			Value: "no",
			Usage: "Should save to database (yes/no))",
		},
	}

	app.Version = "1.0"

	app.Action = func(c *cli.Context) error {
		var args []string
		if c.NArg() > 0 {
			args = c.Args()
			personName := args[0]
			marks := args[1:len(args)]
			log.Println("Person:", personName)
			log.Println("marks", marks)
		}

		if c.String("save") == "no" {
			log.Println("Skipping saving to the database")
		} else {
			log.Println("Saving to the database", args)
		}
		return nil
	}

	app.Run(os.Args)
}