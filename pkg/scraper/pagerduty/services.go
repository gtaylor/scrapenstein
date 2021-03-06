package pagerduty

import (
	"context"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jackc/pgx/v4"
	"time"
)

type ScrapeServicesOptions struct{}

// Scrape and store Pagerduty Services.
func ScrapeServices(dbConn *pgx.Conn, pdClient *pagerduty.Client, options ScrapeServicesOptions) (int, error) {
	listOptions := pagerduty.ListServiceOptions{
		APIListObject: defaultAPIListObject(),
	}
	totalTeams := 0
	for {
		response, err := pdClient.ListServices(listOptions)
		if err != nil {
			return totalTeams, err
		}
		for _, service := range response.Services {
			if err := storeService(dbConn, &service); err != nil {
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
	) VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8
	) ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			name=excluded.name,
			description=excluded.description,
			created_at=excluded.created_at,
			last_incident=excluded.last_incident,
			escalation_policy_id=excluded.escalation_policy_id,
			team_ids=excluded.team_ids
`

func storeService(dbConn *pgx.Conn, service *pagerduty.Service) error {
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
	_, err = dbConn.Exec(
		context.Background(), storeServiceQuery,
		service.ID, service.Summary, service.Name, service.Description, createdAt, lastIncidentPt,
		service.EscalationPolicy.ID, teamIds)
	return err
}
