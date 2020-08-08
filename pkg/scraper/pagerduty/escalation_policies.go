package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

// Scrape and store Pagerduty Escalation Policies.
func ScrapeEscalationPolicies(dbUrl string, authToken string) (int, error) {
	listOptions := pagerduty.ListEscalationPoliciesOptions{
		APIListObject: defaultAPIListObject(),
	}
	client := pagerduty.NewClient(authToken)
	totalPolicies := 0
	for {
		response, err := client.ListEscalationPolicies(listOptions)
		if err != nil {
			return totalPolicies, err
		}
		for _, escalationPolicy := range response.EscalationPolicies {
			if err := storeEscalationPolicy(dbUrl, &escalationPolicy); err != nil {
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
		VALUES($1, $2, $3)
	ON CONFLICT (id)
		DO UPDATE SET name = excluded.name, description = excluded.description`

func storeEscalationPolicy(dbUrl string, escalationPolicy *pagerduty.EscalationPolicy) error {
	_, err := db.SingleExec(
		dbUrl, storeEscalationPolicyQuery,
		escalationPolicy.ID, escalationPolicy.Name, escalationPolicy.Description)
	return err
}
