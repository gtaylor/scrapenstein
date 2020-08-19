package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

type ScrapePrioritiesOptions struct{}

// Scrape and store Pagerduty Priorities.
func ScrapePriorities(dbOptions db.DatabaseOptions, pdOptions PagerDutyOptions, options ScrapePrioritiesOptions) (int, error) {
	client := newPDClient(pdOptions)
	response, err := client.ListPriorities()
	if err != nil {
		return 0, err
	}
	totalPriorities := 0
	for _, priority := range response.Priorities {
		if err := storePriority(dbOptions, &priority); err != nil {
			return totalPriorities, err
		}
		totalPriorities += 1
	}
	return totalPriorities, nil
}

const storePriorityQuery = `
	INSERT INTO pagerduty_priorities (
		id, summary, name, description
	) VALUES(
		$1, $2, $3, $4
	) ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			name=excluded.name,
			description=excluded.description
`

func storePriority(dbOptions db.DatabaseOptions, priority *pagerduty.PriorityProperty) error {
	_, err := db.SingleExec(
		dbOptions, storePriorityQuery,
		priority.ID, priority.Summary, priority.Name, priority.Description)
	return err
}
