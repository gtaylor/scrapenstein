// package pagerduty contains a set of scrapers for PagerDuty.
// See also: https://developer.pagerduty.com/api-reference/
package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"time"
)

// Max paginatable objects
// See: https://developer.pagerduty.com/docs/rest-api-v2/pagination/
const maxListOffset = 10000

// Parses the standard ISO 8601 DateTimes that the PD API returns.
// See: https://developer.pagerduty.com/docs/rest-api-v2/types/
func parseDateTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// Convenience wrapper to provide consistent pagination behavior.
func continuePaginating(listObj pagerduty.APIListObject, total int) bool {
	if listObj.More != true {
		return false
	}
	if total >= maxListOffset {
		return false
	}
	return true
}

func defaultAPIListObject() pagerduty.APIListObject {
	return pagerduty.APIListObject{Offset: 0, Limit: 100}
}
