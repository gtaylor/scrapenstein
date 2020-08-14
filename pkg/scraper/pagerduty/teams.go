package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

type ScrapeTeamsOptions struct{}

// Scrape and store Pagerduty Teams.
func ScrapeTeams(dbOptions db.DatabaseOptions, pdOptions PagerDutyOptions, options ScrapeTeamsOptions) (int, error) {
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
			if err := storeTeam(dbOptions, &team); err != nil {
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
	) VALUES($1, $2, $3, $4)
	ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			name=excluded.name,
			description=excluded.description
`

func storeTeam(dbOptions db.DatabaseOptions, team *pagerduty.Team) error {
	_, err := db.SingleExec(
		dbOptions, storeTeamQuery,
		team.ID, team.Summary, team.Name, team.Description)
	return err
}
