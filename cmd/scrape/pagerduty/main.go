package pagerduty

import (
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/gtaylor/scrapenstein/pkg/db"
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
		Flags: pagerDutyFlags(),
		Before: func(c *cli.Context) error {
			return pagerDutyValidators(c)
		},
		Action: func(c *cli.Context) error {
			pdOptions := pagerDutyOptionsFromCtx(c)
			dbOptions := common.DatabaseOptionsFromCtx(c)
			dbConn, err := db.Connect(dbOptions)
			if err != nil {
				return err
			}
			options := pagerduty.ScrapePrioritiesOptions{}
			logrus.Info("Beginning scrape of PagerDuty Priorities.")
			numScraped, err := pagerduty.ScrapePriorities(dbConn, pdOptions, options)
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
		Flags: pagerDutyFlags(),
		Before: func(c *cli.Context) error {
			return pagerDutyValidators(c)
		},
		Action: func(c *cli.Context) error {
			pdOptions := pagerDutyOptionsFromCtx(c)
			dbOptions := common.DatabaseOptionsFromCtx(c)
			dbConn, err := db.Connect(dbOptions)
			if err != nil {
				return err
			}
			options := pagerduty.ScrapeTeamsOptions{}
			logrus.Info("Beginning scrape of PagerDuty Teams.")
			numScraped, err := pagerduty.ScrapeTeams(dbConn, pdOptions, options)
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
		Flags: pagerDutyFlags(),
		Before: func(c *cli.Context) error {
			return pagerDutyValidators(c)
		},
		Action: func(c *cli.Context) error {
			pdOptions := pagerDutyOptionsFromCtx(c)
			dbOptions := common.DatabaseOptionsFromCtx(c)
			dbConn, err := db.Connect(dbOptions)
			if err != nil {
				return err
			}
			options := pagerduty.ScrapeEscalationPoliciesOptions{}
			logrus.Info("Beginning scrape of PagerDuty Escalation Policies.")
			numScraped, err := pagerduty.ScrapeEscalationPolicies(dbConn, pdOptions, options)
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
		Flags: pagerDutyFlags(),
		Before: func(c *cli.Context) error {
			return pagerDutyValidators(c)
		},
		Action: func(c *cli.Context) error {
			pdOptions := pagerDutyOptionsFromCtx(c)
			dbOptions := common.DatabaseOptionsFromCtx(c)
			dbConn, err := db.Connect(dbOptions)
			if err != nil {
				return err
			}
			options := pagerduty.ScrapeServicesOptions{}
			logrus.Info("Beginning scrape of PagerDuty Services.")
			numScraped, err := pagerduty.ScrapeServices(dbConn, pdOptions, options)
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
		Flags: pagerDutyFlags(
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
			&cli.StringSliceFlag{
				Name:  "team-id",
				Usage: "Only scrape incidents for the specified teams.",
			},
			&cli.StringSliceFlag{
				Name:  "service-id",
				Usage: "Only scrape incidents for the specified service.",
			},
		),
		Before: func(c *cli.Context) error {
			return pagerDutyValidators(c)
		},
		Action: func(c *cli.Context) error {
			pdOptions := pagerDutyOptionsFromCtx(c)
			dbOptions := common.DatabaseOptionsFromCtx(c)
			dbConn, err := db.Connect(dbOptions)
			if err != nil {
				return err
			}

			teamIds := c.StringSlice("team-id")
			if len(teamIds) > 0 {
				logrus.Infof("Limiting scrape to team IDs: %v", teamIds)
			}
			serviceIds := c.StringSlice("service-id")
			if len(serviceIds) > 0 {
				logrus.Infof("Limiting scrape to service IDs: %v", serviceIds)
			}

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

			options := pagerduty.ScrapeIncidentsOptions{
				SinceTime:  sinceTime,
				UntilTime:  untilTime,
				TeamIds:    teamIds,
				ServiceIds: serviceIds,
			}

			logrus.Infof("Beginning scrape of PagerDuty Incidents between %s and %s.",
				sinceTime.Format(time.RFC3339), untilTime.Format(time.RFC3339))
			numScraped, err := pagerduty.ScrapeIncidents(dbConn, pdOptions, options)
			if err != nil {
				return err
			}
			logrus.Infof("Successfully scraped %d PagerDuty Incidents.", numScraped)
			return nil
		},
	}
}
