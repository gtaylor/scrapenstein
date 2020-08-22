package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeOrganizationsOptions struct{}

func ScrapeOrganizations(dbConn *pgx.Conn, ghClient *github.Client, options ScrapeOrganizationsOptions) (int, error) {
	listAllOpts := github.OrganizationsListOptions{}
	totalOrgs := 0
	for {
		orgs, response, err := ghClient.Organizations.ListAll(context.Background(), &listAllOpts)
		if err != nil {
			return totalOrgs, err
		}
		for _, org := range orgs {
			if err := storeOrganization(dbConn, org); err != nil {
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
	INSERT INTO github_organizations (
		id, login, avatar_url
	) VALUES(
		$1, $2, $3
	) ON CONFLICT (id)
		DO UPDATE SET 
			id=excluded.id,
			login=excluded.login,
			avatar_url=excluded.avatar_url`

func storeOrganization(dbConn *pgx.Conn, org *github.Organization) error {
	_, err := dbConn.Exec(
		context.Background(), storeOrganizationQuery,
		org.ID, org.Login, org.AvatarURL)
	return err
}
