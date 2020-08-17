package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

type ScrapeOrganizationsOptions struct{}

func ScrapeOrganizations(dbOptions db.DatabaseOptions, ghOptions GitHubOptions, options ScrapeOrganizationsOptions) (int, error) {
	client, err := newGHClient(ghOptions)
	if err != nil {
		return 0, err
	}
	listAllOpts := github.OrganizationsListOptions{}
	totalOrgs := 0
	for {
		orgs, response, err := client.Organizations.ListAll(context.Background(), &listAllOpts)
		if err != nil {
			return totalOrgs, err
		}
		for _, org := range orgs {
			if err := storeOrganization(dbOptions, org); err != nil {
				return totalOrgs, err
			}
			totalOrgs += 1
		}
		if !continuePaginating(response) {
			break
		}
		listAllOpts.Page = response.NextPage
	}
	return totalOrgs, nil
}

const storeOrganizationQuery = `
	INSERT INTO github_organizations (id, login, url, avatar_url)
		VALUES($1, $2, $3, $4)
	ON CONFLICT (id)
		DO UPDATE SET 
			id=excluded.id, 
			login=excluded.login, 
			url=excluded.url,
			avatar_url=excluded.avatar_url`

func storeOrganization(dbOptions db.DatabaseOptions, org *github.Organization) error {
	_, err := db.SingleExec(
		dbOptions, storeOrganizationQuery,
		org.ID, org.Login, org.URL, org.AvatarURL)
	return err
}
