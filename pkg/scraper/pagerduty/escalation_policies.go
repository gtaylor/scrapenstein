package pagerduty

import (
	"context"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jackc/pgx/v4"
)

type ScrapeEscalationPoliciesOptions struct{}

// Scrape and store Pagerduty Escalation Policies.
func ScrapeEscalationPolicies(dbConn *pgx.Conn, pdOptions PagerDutyOptions, options ScrapeEscalationPoliciesOptions) (int, error) {
	listOptions := pagerduty.ListEscalationPoliciesOptions{
		APIListObject: defaultAPIListObject(),
	}
	client := newPDClient(pdOptions)
	totalPolicies := 0
	for {
		response, err := client.ListEscalationPolicies(listOptions)
		if err != nil {
			return totalPolicies, err
		}
		for _, escalationPolicy := range response.EscalationPolicies {
			if err := storeEscalationPolicy(dbConn, &escalationPolicy); err != nil {
				return totalPolicies, err
			}
			totalPolicies += 1
		}
		if !continuePaginating(response.APIListObject, totalPolicies) {
			break
		}
		listOptions.Offset = uint(totalPolicies)
	}
	return totalPolicies, nil
}

const storeEscalationPolicyQuery = `
	INSERT INTO pagerduty_escalation_policies (id, name, description)
	VALUES(
		$1, $2, $3
	) ON CONFLICT (id)
		DO UPDATE SET 
			name=excluded.name, 
			description=excluded.description`

func storeEscalationPolicy(dbConn *pgx.Conn, escalationPolicy *pagerduty.EscalationPolicy) error {
	_, err := dbConn.Exec(
		context.Background(), storeEscalationPolicyQuery,
		escalationPolicy.ID, escalationPolicy.Name, escalationPolicy.Description)
	return err
}
