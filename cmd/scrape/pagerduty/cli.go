package pagerduty

import (
	"errors"
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/urfave/cli/v2"
)

const authTokenFlagName = "auth-token"

func authTokenFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    authTokenFlagName,
		Usage:   "PagerDuty auth token",
		EnvVars: []string{"PAGERDUTY_AUTH_TOKEN"},
	}
}

// Most or all of the PagerDuty subcommands have the same flags.
func pdFlags(others ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		common.DatabaseURLFlag(),
		authTokenFlag(),
	}, others...)
}

func mustSetAuthToken(c *cli.Context) error {
	if c.String(authTokenFlagName) == "" {
		return errors.New("PagerDuty auth token must be set via --auth-token or PAGERDUTY_AUTH_TOKEN.")
	}
	return nil
}

// The pre-command (Before) validation is the same for all PD subcommands as well.
func pdValidators(c *cli.Context) error {
	if err := common.MustSetDatabaseURL(c); err != nil {
		return err
	}
	if err := mustSetAuthToken(c); err != nil {
		return err
	}
	return nil
}
