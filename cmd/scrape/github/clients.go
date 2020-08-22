package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/gtaylor/scrapenstein/pkg/db"
	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// Holds GitHub client options.
type gitHubOptions struct {
	AccessToken string
	BaseURL     string
	UploadURL   string
}

// Pulls GH options from the CLI and returns a gitHubOptions instance.
func githubOptionsFromCtx(c *cli.Context) gitHubOptions {
	return gitHubOptions{
		AccessToken: c.String(accessTokenFlagName),
		BaseURL:     c.String(baseURLFlagName),
		UploadURL:   c.String(uploadURLFlagName),
	}
}

func newGHClient(options gitHubOptions) (*github.Client, error) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: options.AccessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	if options.BaseURL != "" && options.UploadURL != "" {
		return github.NewEnterpriseClient(options.BaseURL, options.UploadURL, httpClient)
	} else {
		return github.NewClient(httpClient), nil
	}
}

func newClients(c *cli.Context) (*pgx.Conn, *github.Client, error) {
	dbOptions := common.DatabaseOptionsFromCtx(c)
	dbConn, err := db.Connect(dbOptions)
	if err != nil {
		return nil, nil, err
	}
	ghOptions := githubOptionsFromCtx(c)
	ghClient, err := newGHClient(ghOptions)
	if err != nil {
		return nil, nil, err
	}
	return dbConn, ghClient, nil
}
