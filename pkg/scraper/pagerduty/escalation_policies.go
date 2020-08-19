package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

type ScrapeEscalationPoliciesOptions struct{}

// Scrape and store Pagerduty Escalation Policies.
func ScrapeEscalationPolicies(dbOptions db.DatabaseOptions, pdOptions PagerDutyOptions, options ScrapeEscalationPoliciesOptions) (int, error) {
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
			if err := storeEscalationPolicy(dbOptions, &escalationPolicy); err != nil {
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

func storeEscalationPolicy(dbOptions db.DatabaseOptions, escalationPolicy *pagerduty.EscalationPolicy) error {
	_, err := db.SingleExec(
		dbOptions, storeEscalationPolicyQuery,
		escalationPolicy.ID, escalationPolicy.Name, escalationPolicy.Description)
	return err
}
