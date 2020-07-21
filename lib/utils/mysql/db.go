package mysql

import (
	"fmt"
	"database/sql"
)

// NewDB returns a new *sql.DB connector with a mysql database. It also pings the database to ensures the connection is working.
func NewDB(user, password, host, port, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}