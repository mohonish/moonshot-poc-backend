package db

import (
	"encoding/json"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	//"github.com/mohonish/moonshot-backend/models"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"log"
)

// const databaseURL string = "postgres://YourUserName:YourPassword@YourHost:5432/YourDatabase"
const databaseURL string = "postgres://postgres:postgres@localhost:5432/iss_data"

// const getLatestSQLQuery string = "select * from location where time > now() - interval '10 minutes' order by time desc limit 1"
const getLatestSQLQuery string = "select * from location order by time desc limit 1"

const remoteDataService string = "http://api.open-notify.org/iss-now.json"

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

func pollRemoteDataService() {
	apiClient := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, remoteDataService, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := apiClient.Do(req)
	if getErr != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	type locationResponse struct {

	}
	parsedResponse := locationResponse{}
	jsonErr := json.Unmarshal(body, &parsedResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// persist to db.

}