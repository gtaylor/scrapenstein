package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

// Scrape and store Pagerduty Teams.
func ScrapeTeams(dbUrl string, authToken string) (int, error) {
	listOptions := pagerduty.ListTeamOptions{
		APIListObject: defaultAPIListObject(),
	}
	client := pagerduty.NewClient(authToken)
	totalTeams := 0
	for {
		response, err := client.ListTeams(listOptions)
		if err != nil {
			return totalTeams, err
		}
		for _, team := range response.Teams {
			if err := storeTeam(dbUrl, &team); err != nil {
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

func storeTeam(dbUrl string, team *pagerduty.Team) error {
	_, err := db.SingleExec(
		dbUrl, storeTeamQuery,
		team.ID, team.Summary, team.Name, team.Description)
	return err
}
