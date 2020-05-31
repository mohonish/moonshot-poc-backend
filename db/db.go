package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

// const databaseURL string = "postgres://YourUserName:YourPassword@YourHost:5432/YourDatabase"
const databaseURL string = "postgres://postgres:postgres@localhost:5432/iss_data"
// const getLatestSQLQuery string = "select * from location where time > now() - interval '10 minutes' order by time desc limit 1"
const getLatestSQLQuery string = "select * from location order by time desc limit 1"

var conn *pgx.Conn

func SetupDBFetch() {
	// Connect to DB
	var err error
	conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Connected to database\n")
	}
	// defer conn.Close(context.Background())
}

func GetLatestLocation() (timestamp time.Time, lat float64, long float64, err error) {
	err = conn.QueryRow(context.Background(), getLatestSQLQuery).Scan(&timestamp, &lat, &long)
	if err != nil {
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database access error: %v\n", err)
	}
	return
}
