// package github contains a set of scrapers for GitHub.
// See also: https://docs.github.com/en/rest
package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"time"
)

// Holds GitHub client options.
type GitHubOptions struct {
	AccessToken string
	BaseURL     string
	UploadURL   string
}

func newGHClient(options GitHubOptions) (*github.Client, error) {
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

// Convenience wrapper to provide consistent pagination behavior.
func continuePaginating(response *github.Response) bool {
	if response.NextPage == 0 {
		return false
	}
	return true
}

func ghTsToTimeOrNull(ghTs *github.Timestamp) *time.Time {
	if ghTs == nil {
		return nil
	}
	return &ghTs.Time
}
