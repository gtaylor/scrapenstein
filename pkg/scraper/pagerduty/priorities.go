package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/pkg/db"
)

type ScrapePrioritiesOptions struct{}

// Scrape and store Pagerduty Priorities.
func ScrapePriorities(dbUrl string, authToken string, options ScrapePrioritiesOptions) (int, error) {
	client := pagerduty.NewClient(authToken)
	totalPriorities := 0
	response, err := client.ListPriorities()
	if err != nil {
		return totalPriorities, err
	}
	for _, priority := range response.Priorities {
		if err := storePriority(dbUrl, &priority); err != nil {
			return totalPriorities, err
		}
		totalPriorities += 1
	}
	return totalPriorities, nil
}

const storePriorityQuery = `
	INSERT INTO pagerduty_priorities (
		id, summary, name, description
	) VALUES($1, $2, $3, $4)
	ON CONFLICT (id)
		DO UPDATE SET 
			summary=excluded.summary, 
			name=excluded.name,
			description=excluded.description
`

func storePriority(dbUrl string, priority *pagerduty.PriorityProperty) error {
	_, err := db.SingleExec(
		dbUrl, storePriorityQuery,
		priority.ID, priority.Summary, priority.Name, priority.Description)
	return err
}
