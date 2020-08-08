package common

import (
	"errors"
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
