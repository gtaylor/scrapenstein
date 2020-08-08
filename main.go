package main

import (
	"github.com/gtaylor/scrapenstein/cmd/scrape"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "scrapenstein",
		Usage: "Scrape all the things",
		Commands: []*cli.Command{
			scrape.Command(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
