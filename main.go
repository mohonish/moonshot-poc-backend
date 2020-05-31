package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/mohonish/moonshot-backend/version"
	"log"
	"net/http"
	"os"
)

const contentTypeKey string = "Content-Type"
const contentTypeVal string = "application/json"
// const databaseURL string = "postgres://YourUserName:YourPassword@YourHost:5432/YourDatabase"
const databaseURL string = "postgres://postgres:postgres@localhost:5432/iss_data"
const port string = ":8081"

func main() {
	// Set CLI Flags
	setCLIFlags()

	// Connect to DB
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Connected to database\n")
	}
	defer conn.Close(context.Background())

	// Create mux base router
	r := mux.NewRouter()

	// Create mux subrouter
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/iss/coordinates", baseRouteGET).Methods(http.MethodGet)
	api.HandleFunc("", routeNotFound)

	// Listen on 8081
	fmt.Fprintf(os.Stdout, "Server listening on 8081\n")
	log.Fatal(http.ListenAndServe(port, r))
}

func baseRouteGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func setCLIFlags() {
	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		return
	}
}
