package pagerduty

import (
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/gtaylor/scrapenstein/pkg/scraper/pagerduty"
	"github.com/karrick/tparse/v2"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"time"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "pagerduty",
		Usage: "PagerDuty scraper",
		Subcommands: []*cli.Command{
			priorityCommand(),
			teamCommand(),
			escalationPolicyCommand(),
			serviceCommand(),
			incidentCommand(),
		},
	}
}

func priorityCommand() *cli.Command {
	return &cli.Command{
		Name:  "priority",
		Usage: "Scrape Priorities",
		Flags: pdFlags(),
		Before: func(c *cli.Context) error {
			return pdValidators(c)
		},
		Action: func(c *cli.Context) error {
			authToken := c.String(authTokenFlagName)
			dbUrl := c.String(common.DatabaseURLFlagName)
			logrus.Info("Beginning scrape of PagerDuty Priorities.")
			numScraped, err := pagerduty.ScrapePriorities(dbUrl, authToken)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d PagerDuty Priorities.", numScraped)
			return nil
		},
	}
}

func teamCommand() *cli.Command {
	return &cli.Command{
		Name:  "team",
		Usage: "Scrape Teams",
		Flags: pdFlags(),
		Before: func(c *cli.Context) error {
			return pdValidators(c)
		},
		Action: func(c *cli.Context) error {
			authToken := c.String(authTokenFlagName)
			dbUrl := c.String(common.DatabaseURLFlagName)
			logrus.Info("Beginning scrape of PagerDuty Teams.")
			numScraped, err := pagerduty.ScrapeTeams(dbUrl, authToken)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d PagerDuty Teams.", numScraped)
			return nil
		},
	}
}

func escalationPolicyCommand() *cli.Command {
	return &cli.Command{
		Name:  "escalation",
		Usage: "Scrape Escalation Policies",
		Flags: pdFlags(),
		Before: func(c *cli.Context) error {
			return pdValidators(c)
		},
		Action: func(c *cli.Context) error {
			authToken := c.String(authTokenFlagName)
			dbUrl := c.String(common.DatabaseURLFlagName)
			logrus.Info("Beginning scrape of PagerDuty Escalation Policies.")
			numScraped, err := pagerduty.ScrapeEscalationPolicies(dbUrl, authToken)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d PagerDuty Escalation Policies.", numScraped)
			return nil
		},
	}
}

func serviceCommand() *cli.Command {
	return &cli.Command{
		Name:  "service",
		Usage: "Scrape Services",
		Flags: pdFlags(),
		Before: func(c *cli.Context) error {
			return pdValidators(c)
		},
		Action: func(c *cli.Context) error {
			authToken := c.String(authTokenFlagName)
			dbUrl := c.String(common.DatabaseURLFlagName)
			logrus.Info("Beginning scrape of PagerDuty Services.")
			numScraped, err := pagerduty.ScrapeServices(dbUrl, authToken)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d PagerDuty Services.", numScraped)
			return nil
		},
	}
}

func incidentCommand() *cli.Command {
	return &cli.Command{
		Name:  "incident",
		Usage: "Scrape Incidents",
		Flags: pdFlags(
			&cli.StringFlag{
				Name:  "since",
				Usage: "The start of the date range to search.",
				Value: "now-30d",
			},
			&cli.StringFlag{
				Name:  "until",
				Usage: "The end of the date range to search.",
				Value: "now",
			},
		),
		Before: func(c *cli.Context) error {
			return pdValidators(c)
		},
		Action: func(c *cli.Context) error {
			authToken := c.String(authTokenFlagName)
			dbUrl := c.String(common.DatabaseURLFlagName)

			since := c.String("since")
			sinceTime, err := tparse.ParseNow(time.RFC3339, since)
			if err != nil {
				return err
			}

			until := c.String("until")
			untilTime, err := tparse.ParseNow(time.RFC3339, until)
			if err != nil {
				return err
			}

			logrus.Infof("Beginning scrape of PagerDuty Incidents between %s and %s.",
				sinceTime.Format(time.RFC3339), untilTime.Format(time.RFC3339))
			numScraped, err := pagerduty.ScrapeIncidents(dbUrl, authToken, sinceTime, untilTime)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d PagerDuty Incidents.", numScraped)
			return nil
		},
	}
}
