package common

import (
	"errors"
	"github.com/gtaylor/scrapenstein/pkg/db"
	"github.com/urfave/cli/v2"
)

const DatabaseURLFlagName = "database-url"

func MustSetDatabaseURL(c *cli.Context) error {
	if c.String(DatabaseURLFlagName) == "" {
		return errors.New("Database URL must be set via --database-url or DATABASE_URL.")
	}
	return nil
}

func DatabaseURLFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    DatabaseURLFlagName,
		Usage:   "Database URL",
		EnvVars: []string{"DATABASE_URL"},
	}
}

// Pulls DB configs from CLI and returns a DatabaseOptions instance.
func DatabaseOptionsFromCtx(c *cli.Context) db.DatabaseOptions {
	return db.DatabaseOptions{
		URL: c.String(DatabaseURLFlagName),
	}
}
