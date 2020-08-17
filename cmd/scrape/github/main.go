package github

import (
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/gtaylor/scrapenstein/pkg/scraper/github"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "github",
		Usage: "GitHub scraper",
		Subcommands: []*cli.Command{
			organizationsCommand(),
			repositoryCommand(),
		},
	}
}

func organizationsCommand() *cli.Command {
	return &cli.Command{
		Name:  "organizations",
		Usage: "Scrape Organizations",
		Flags: gitHubFlags(),
		Before: func(c *cli.Context) error {
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbOptions := common.DatabaseOptionsFromCtx(c)
			ghOptions := githubOptionsFromCtx(c)
			options := github.ScrapeOrganizationsOptions{}
			logrus.Info("Beginning scrape of GitHub Organizations.")
			numScraped, err := github.ScrapeOrganizations(dbOptions, ghOptions, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d GitHub Organizations.", numScraped)
			return nil
		},
	}
}

func repositoryCommand() *cli.Command {
	return &cli.Command{
		Name:      "repository",
		Usage:     "Scrape a Repository",
		Flags:     gitHubFlags(),
		ArgsUsage: "<owner> <repo>",
		Before: func(c *cli.Context) error {
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbOptions := common.DatabaseOptionsFromCtx(c)
			ghOptions := githubOptionsFromCtx(c)
			options := github.ScrapeRepositoryOptions{
				Owner: c.Args().Get(0),
				Repo:  c.Args().Get(1),
			}
			logrus.Infof("Beginning scrape GitHub Repository: %s/%s", options.Owner, options.Repo)
			err := github.ScrapeRepository(dbOptions, ghOptions, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped Github Repository %s/%s", options.Owner, options.Repo)
			return nil
		},
	}
}
