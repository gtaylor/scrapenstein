package pagerduty

import (
	"context"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jackc/pgx/v4"
)

type ScrapeTeamsOptions struct{}

// Scrape and store Pagerduty Teams.
func ScrapeTeams(dbConn *pgx.Conn, pdOptions PagerDutyOptions, options ScrapeTeamsOptions) (int, error) {
	listOptions := pagerduty.ListTeamOptions{
		APIListObject: defaultAPIListObject(),
	}
	client := newPDClient(pdOptions)
	totalTeams := 0
	for {
		response, err := client.ListTeams(listOptions)
		if err != nil {
			return totalTeams, err
		}
		for _, team := range response.Teams {
			if err := storeTeam(dbConn, &team); err != nil {
				return totalTeams, err
			}
			totalTeams += 1
		}
		if !continuePaginating(response.APIListObject, totalTeams) {
			break
		}
		listOptions.Offset = uint(totalTeams)
	}
	return totalTeams, nil
}

const storeTeamQuery = `
	INSERT INTO pagerduty_teams (
		id, summary, name, description
	) VALUES(
		$1, $2, $3, $4
	) ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			name=excluded.name,
			description=excluded.description
`

func storeTeam(dbConn *pgx.Conn, team *pagerduty.Team) error {
	_, err := dbConn.Exec(
		context.Background(), storeTeamQuery,
		team.ID, team.Summary, team.Name, team.Description)
	return err
}
