package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/gtaylor/scrapenstein/cmd/scrape/common"
	"github.com/gtaylor/scrapenstein/pkg/db"
	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

// Holds PagerDuty client options.
type pagerDutyOptions struct {
	AuthToken string
}

func pagerDutyOptionsFromCtx(c *cli.Context) pagerDutyOptions {
	return pagerDutyOptions{
		AuthToken: c.String(authTokenFlagName),
	}
}

func newPDClient(pdOptions pagerDutyOptions) *pagerduty.Client {
	return pagerduty.NewClient(pdOptions.AuthToken)
}

func newClients(c *cli.Context) (*pgx.Conn, *pagerduty.Client, error) {
	dbOptions := common.DatabaseOptionsFromCtx(c)
	dbConn, err := db.Connect(dbOptions)
	if err != nil {
		return nil, nil, err
	}
	pdOptions := pagerDutyOptionsFromCtx(c)
	pdClient := newPDClient(pdOptions)
	return dbConn, pdClient, nil
}
