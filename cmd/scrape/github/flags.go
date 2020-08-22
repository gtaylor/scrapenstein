package github

import (
	"errors"
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/urfave/cli/v2"
)

const accessTokenFlagName = "access-token"
const baseURLFlagName = "base-url"
const uploadURLFlagName = "upload-url"

func accessTokenFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    accessTokenFlagName,
		Usage:   "GitHub access token",
		EnvVars: []string{"GITHUB_ACCESS_TOKEN"},
	}
}

func baseURLFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    baseURLFlagName,
		Usage:   "GitHub base URL (GHE)",
		EnvVars: []string{"GITHUB_BASE_URL"},
	}
}

func uploadURLFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    uploadURLFlagName,
		Usage:   "GitHub upload URL (GHE)",
		EnvVars: []string{"GITHUB_UPLOAD_URL"},
	}
}

// Most or all of the GitHub subcommands have the same flags.
func gitHubFlags(others ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		common.DatabaseURLFlag(),
		accessTokenFlag(),
		baseURLFlag(),
		uploadURLFlag(),
	}, others...)
}

func mustSetAccessToken(c *cli.Context) error {
	if c.String(accessTokenFlagName) == "" {
		return errors.New("GitHub access token must be set via --access-token or GITHUB_ACCESS_TOKEN.")
	}
	if c.String(baseURLFlagName) != "" && c.String(uploadURLFlagName) == "" {
		return errors.New("A GitHub base URL was provided without a corresponding upload URL.")
	}
	if c.String(baseURLFlagName) == "" && c.String(uploadURLFlagName) != "" {
		return errors.New("A GitHub upload URL was provided without a corresponding base URL.")
	}
	return nil
}

// The pre-command (Before) validation is the same for all GitHub subcommands as well.
func gitHubValidators(c *cli.Context) error {
	if err := common.MustSetDatabaseURL(c); err != nil {
		return err
	}
	if err := mustSetAccessToken(c); err != nil {
		return err
	}
	return nil
}

// Verify that an org and repo name were passed in.
// NOTE: org and repo must be args 1 and 2 respectively.
func orgAndRepoValidator(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("You must pass in the org and repo name.")
	}
	if c.Args().Len() < 2 {
		return errors.New("You must pass in the repo name.")
	}
	return nil
}

func orgValidator(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("You must pass in the org name.")
	}
	return nil
}
