package github

import (
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
			usersCommand(),
			teamsCommand(),
			repositoryCommand(),
			commitsCommand(),
			pullRequestsCommand(),
			issuesCommand(),
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
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapeOrganizationsOptions{}
			logrus.Info("Beginning scrape of GitHub Organizations.")
			numScraped, err := github.ScrapeOrganizations(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d GitHub Organizations.", numScraped)
			return nil
		},
	}
}

func usersCommand() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Scrape all Users on your GH instance",
		Flags: gitHubFlags(),
		Before: func(c *cli.Context) error {
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapeUsersOptions{}
			logrus.Info("Beginning scrape of all Users.")
			numScraped, err := github.ScrapeUsers(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d users.", numScraped)
			return nil
		},
	}
}

func teamsCommand() *cli.Command {
	return &cli.Command{
		Name:      "teams",
		Usage:     "Scrape an Organization's Teams",
		Flags:     gitHubFlags(),
		ArgsUsage: "<org>",
		Before: func(c *cli.Context) error {
			if err := orgValidator(c); err != nil {
				return err
			}
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapeTeamsOptions{
				Org: c.Args().Get(0),
			}
			logrus.Infof("Beginning scrape of GitHub Teams from org %s.", options.Org)
			numScraped, err := github.ScrapeTeams(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d GitHub Teams from org %s.", numScraped, options.Org)
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
			if err := orgAndRepoValidator(c); err != nil {
				return err
			}
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapeRepositoryOptions{
				OrgRepoAndRepoId: github.OrgRepoAndRepoId{
					Owner: c.Args().Get(0),
					Repo:  c.Args().Get(1),
				},
			}
			logrus.Infof("Beginning scrape of GitHub Repository %s/%s.", options.Owner, options.Repo)
			err = github.ScrapeRepository(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped Github Repository %s/%s.", options.Owner, options.Repo)
			return nil
		},
	}
}

func commitsCommand() *cli.Command {
	return &cli.Command{
		Name:  "commits",
		Usage: "Scrape a Repository's Commits",
		Flags: gitHubFlags(
			&cli.BoolFlag{
				Name:  "scrape-stats",
				Usage: "Scrape commit stats summary (slow).",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "scrape-files",
				Usage: "Scrape commit file changes (very slow).",
				Value: false,
			},
		),
		ArgsUsage: "<owner> <repo>",
		Before: func(c *cli.Context) error {
			if err := orgAndRepoValidator(c); err != nil {
				return err
			}
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapeCommitsOptions{
				OrgRepoAndRepoId: github.OrgRepoAndRepoId{
					Owner: c.Args().Get(0),
					Repo:  c.Args().Get(1),
				},
				ScrapeStats: c.Bool("scrape-stats"),
				ScrapeFiles: c.Bool("scrape-files"),
			}
			if options.ScrapeStats {
				logrus.Infof("Enabling the scraping of commit change stats.")
			}
			if options.ScrapeFiles {
				logrus.Infof("Enabling the scraping of per-file change stats.")
			}
			logrus.Infof("Beginning scrape of GitHub Commits from %s/%s.", options.Owner, options.Repo)
			numScraped, err := github.ScrapeCommits(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d Github Commits from %s/%s.", numScraped, options.Owner, options.Repo)
			return nil
		},
	}
}

func pullRequestsCommand() *cli.Command {
	return &cli.Command{
		Name:  "pullrequests",
		Usage: "Scrape a Repository's Pull Requests",
		Flags: gitHubFlags(
			&cli.BoolFlag{
				Name:  "scrape-stats",
				Usage: "Scrape PR stats summary (very slow).",
				Value: false,
			},
		),
		ArgsUsage: "<owner> <repo>",
		Before: func(c *cli.Context) error {
			if err := orgAndRepoValidator(c); err != nil {
				return err
			}
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapePullRequestsOptions{
				OrgRepoAndRepoId: github.OrgRepoAndRepoId{
					Owner: c.Args().Get(0),
					Repo:  c.Args().Get(1),
				},
				ScrapeStats: c.Bool("scrape-stats"),
			}
			if options.ScrapeStats {
				logrus.Infof("Enabling the scraping of PR stats.")
			}
			logrus.Infof("Beginning scrape of GitHub Pull Requests from %s/%s.", options.Owner, options.Repo)
			numScraped, err := github.ScrapePullRequests(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d Github Pull Requests from %s/%s.", numScraped, options.Owner, options.Repo)
			return nil
		},
	}
}

func issuesCommand() *cli.Command {
	return &cli.Command{
		Name:      "issues",
		Usage:     "Scrape a Repository's Issues",
		Flags:     gitHubFlags(),
		ArgsUsage: "<owner> <repo>",
		Before: func(c *cli.Context) error {
			if err := orgAndRepoValidator(c); err != nil {
				return err
			}
			return gitHubValidators(c)
		},
		Action: func(c *cli.Context) error {
			dbConn, ghClient, err := newClients(c)
			if err != nil {
				return err
			}
			options := github.ScrapeIssuesOptions{
				OrgRepoAndRepoId: github.OrgRepoAndRepoId{
					Owner: c.Args().Get(0),
					Repo:  c.Args().Get(1),
				},
			}
			logrus.Infof("Beginning scrape of GitHub Issues from %s/%s.", options.Owner, options.Repo)
			numScraped, err := github.ScrapeIssues(dbConn, ghClient, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d Github Issues from %s/%s.", numScraped, options.Owner, options.Repo)
			return nil
		},
	}
}
