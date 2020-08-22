// package github contains a set of scrapers for GitHub.
// See also: https://docs.github.com/en/rest
package github

import (
	"context"
	"errors"
	"github.com/google/go-github/v32/github"
	"time"
)

// Convenience struct used to hold repo identifiers, typically in scrape options structs.
type OrgRepoAndRepoId struct {
	Owner string
	Repo  string
	// Providing RepoID allows for skipping a query for the repo's ID for List calls.
	RepoId int64
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

// If a RepoId is not set, query one from org name + repo name and set it on the struct.
// We use Repo IDs instead of org name + repo name as our primary keys in order to
// gracefully handle repo renames and moves between orgs and owners.
func ensureRepoIdFromAPI(client *github.Client, repoOpts *OrgRepoAndRepoId) error {
	if repoOpts.RepoId > 0 {
		return nil
	}
	repo, _, err := client.Repositories.Get(context.Background(), repoOpts.Owner, repoOpts.Repo)
	if err != nil {
		return err
	}
	repoOpts.RepoId = repo.GetID()
	if repoOpts.RepoId == 0 {
		return errors.New("unable to retrieve a non-zero Repo ID from the API")
	}
	return nil
}
