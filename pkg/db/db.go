package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
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

// Convenience function for obtaining a DB conn and executing a single query.
// Make sure to avoid calling this multiple times in one scraper. If you're going
// to be making multiple queries, call Connect() + Query/Execute() separately.
func SingleExec(dbOptions DatabaseOptions, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	dbConn, err := Connect(dbOptions)
	if err != nil {
		return nil, err
	}
	_, err = dbConn.Exec(context.Background(), sql, arguments...)
	return nil, err
}
