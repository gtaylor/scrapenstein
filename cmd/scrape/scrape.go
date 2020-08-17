package scrape

import (
	"github.com/gtaylor/scrapenstein/cmd/scrape/github"
	"github.com/gtaylor/scrapenstein/cmd/scrape/pagerduty"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "scrape",
		Usage: "Various and sundry scrapers",
		Subcommands: []*cli.Command{
			pagerduty.Command(),
			github.Command(),
		},
	}
}
