package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
	"time"
)

// Scrape and store Pagerduty Incidents.
func ScrapeIncidents(dbUrl string, authToken string, sinceTime time.Time, untilTime time.Time) (int, error) {
	listOptions := pagerduty.ListIncidentsOptions{
		Since:         sinceTime.Format(time.RFC3339),
		Until:         untilTime.Format(time.RFC3339),
		APIListObject: defaultAPIListObject(),
	}
	client := pagerduty.NewClient(authToken)
	totalIncidents := 0
	for {
		response, err := client.ListIncidents(listOptions)
		if err != nil {
			return totalIncidents, err
		}
		for _, incident := range response.Incidents {
			if err := storeIncident(dbUrl, &incident); err != nil {
				return totalIncidents, err
			}
			totalIncidents += 1
		}
		if !continuePaginating(response.APIListObject, totalIncidents) {
			break
		}
		listOptions.Offset = uint(totalIncidents)
	}
	return totalIncidents, nil
}

const storeIncidentQuery = `
	INSERT INTO pagerduty_incidents (
		id, summary, incident_number, created_at, status, title, incident_key, service_id,
		last_status_change_at, escalation_policy_id, team_ids, priority_id, urgency
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			incident_number=excluded.incident_number,
			created_at=excluded.created_at,
			status=excluded.status,
			title=excluded.title,
			incident_key=excluded.incident_key,
			service_id=excluded.service_id,
			last_status_change_at=excluded.last_status_change_at,
			escalation_policy_id=excluded.escalation_policy_id,
			team_ids=excluded.team_ids,
			priority_id=excluded.priority_id,
			urgency=excluded.urgency`

func storeIncident(dbUrl string, incident *pagerduty.Incident) error {
	createdAt, err := parseDateTime(incident.CreatedAt)
	if err != nil {
		return err
	}
	lastStatusChangeAt, err := parseDateTime(incident.LastStatusChangeAt)
	if err != nil {
		return err
	}
	priorityId := ""
	if incident.Priority != nil {
		priorityId = incident.Priority.ID
	}
	teamIds := make([]string, 0)
	if len(incident.Teams) > 0 {
		for _, team := range incident.Teams {
			teamIds = append(teamIds, team.ID)
		}
	}

	_, err = db.SingleExec(
		dbUrl, storeIncidentQuery,
		incident.Id, incident.Summary, incident.IncidentNumber, createdAt,
		incident.Status, incident.Title, incident.IncidentKey, incident.Service.ID,
		lastStatusChangeAt, incident.EscalationPolicy.ID, teamIds, priorityId,
		incident.Urgency)
	return err
}