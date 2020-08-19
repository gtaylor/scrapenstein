package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type DatabaseOptions struct {
	URL string
}

// Connects to the a Postgres DB to be used for storing scraper results.
func Connect(dbOptions DatabaseOptions) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbOptions.URL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	return conn, nil
}
