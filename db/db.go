package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"log"
)

// const databaseURL string = "postgres://YourUserName:YourPassword@YourHost:5432/YourDatabase"
const databaseURL string = "postgres://postgres:postgres@localhost:5432/iss_data"

// const getLatestSQLQuery string = "select * from location where time > now() - interval '10 minutes' order by time desc limit 1"
const getLatestSQLQuery string = "select * from location order by time desc limit 1"
const insertQuery string = `INSERT INTO location(time, latitude, longitude) VALUES ($1, $2, $3)`

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

	schedule(pollRemoteDataService, 5*time.Second)
	//pollRemoteDataService()
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

/* Data Fetch Methods */

func schedule(f func(), interval time.Duration) *time.Ticker {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            f()
        }
    }()
    return ticker
}

func pollRemoteDataService() {
	apiClient := http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest(http.MethodGet, remoteDataService, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := apiClient.Do(req)
	if getErr != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fmt.Fprintf(os.Stdout, "Body: %s", body)

	/*
	var pos struct {
		coordinates	models.Coordinate		`json:"iss_position"`
		timestamp	uint64			`json:"timestamp"`
	}
	jsonErr := json.Unmarshal(body, &pos)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	fmt.Fprintf(os.Stdout, "\n%+v\n", pos)
	*/
	
	result := string(body[:])
	latResult := gjson.Get(result, "iss_position.latitude")
	longResult := gjson.Get(result, "iss_position.longitude")
	// fmt.Fprintf(os.Stdout, "\n%+v %+v\n", lat_result, long_result)
	
	lat, _ := strconv.ParseFloat(latResult.Str, 8)
	long, _ := strconv.ParseFloat(longResult.Str, 8)
	
	// fmt.Fprintf(os.Stdout, "\n%f %f\n", lat, long)

	if lat == 0 || long == 0 {
		log.Fatal("Nil response after parsing")
	}

	// persist to db.
	dbErr := insertToDB(lat, long)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	
}

func insertToDB(lat float64, long float64) error {
	fmt.Fprintf(os.Stdout, "\nInsert -> lat:%f long:%f\n", lat, long)
	_, err := conn.Exec(context.Background(), insertQuery, time.Now(), lat, long)
	return err
}