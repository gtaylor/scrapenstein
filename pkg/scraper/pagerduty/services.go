package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
	"time"
)

// Scrape and store Pagerduty Services.
func ScrapeServices(dbUrl string, authToken string) (int, error) {
	listOptions := pagerduty.ListServiceOptions{
		APIListObject: defaultAPIListObject(),
	}
	client := pagerduty.NewClient(authToken)
	totalTeams := 0
	for {
		response, err := client.ListServices(listOptions)
		if err != nil {
			return totalTeams, err
		}
		for _, service := range response.Services {
			if err := storeService(dbUrl, &service); err != nil {
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

const storeServiceQuery = `
	INSERT INTO pagerduty_services (
		id, summary, name, description, created_at, last_incident, escalation_policy_id, team_ids
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			name=excluded.name,
			description=excluded.description,
			created_at=excluded.created_at,
			last_incident=excluded.last_incident,
			escalation_policy_id=excluded.escalation_policy_id,
			team_ids=excluded.team_ids
`

func storeService(dbUrl string, service *pagerduty.Service) error {
	createdAt, err := parseDateTime(service.CreateAt)
	if err != nil {
		return err
	}
	var lastIncidentPt *time.Time
	// Not all services have had an incident... yet.
	if service.LastIncidentTimestamp != "" {
		lastIncident, err := parseDateTime(service.LastIncidentTimestamp)
		if err != nil {
			return err
		}
		lastIncidentPt = &lastIncident
	}
	teamIds := make([]string, 0)
	if len(service.Teams) > 0 {
		for _, team := range service.Teams {
			teamIds = append(teamIds, team.ID)
		}
	}
	_, err = db.SingleExec(
		dbUrl, storeServiceQuery,
		service.ID, service.Summary, service.Name, service.Description, createdAt, lastIncidentPt,
		service.EscalationPolicy.ID, teamIds)
	return err
}
