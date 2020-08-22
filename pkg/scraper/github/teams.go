package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeTeamsOptions struct {
	Org string
}

func ScrapeTeams(dbConn *pgx.Conn, ghClient *github.Client, options ScrapeTeamsOptions) (int, error) {
	listOpts := github.ListOptions{}
	totalTeams := 0
	for {
		teams, response, err := ghClient.Teams.ListTeams(context.Background(), options.Org, &listOpts)
		if err != nil {
			return totalTeams, err
		}
		for _, team := range teams {
			if err := storeTeam(dbConn, team); err != nil {
				return totalTeams, err
			}
			totalTeams += 1
		}
		if !continuePaginating(response) {
			break
		}
		listOpts.Page = response.NextPage
	}
	return totalTeams, nil
}

const storeTeamQuery = `
	INSERT INTO github_teams (
		id, name, slug, description, privacy, permission, parent_id
	) VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8
	) ON CONFLICT (id)
		DO UPDATE SET
			id=excluded.id,
			name=excluded.name,
			slug=excluded.slug,
			description=excluded.description,
			privacy=excluded.privacy,
			permission=excluded.permission,
			parent_id=excluded.parent_id`

func storeTeam(dbConn *pgx.Conn, team *github.Team) error {
	parent := team.GetParent()
	var parentId *int64
	if parent != nil {
		parentId = parent.ID
	}
	_, err := dbConn.Exec(
		context.Background(), storeTeamQuery,
		team.GetID(), team.GetName(), team.GetSlug(), team.GetDescription(),
		team.GetPrivacy(), team.GetPermission(), parentId)
	return err
}
