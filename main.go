package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mohonish/moonshot-backend/db"
	"github.com/mohonish/moonshot-backend/models"
	"github.com/mohonish/moonshot-backend/version"
	"log"
	"net/http"
	"os"
)

const contentTypeKey string = "Content-Type"
const contentTypeVal string = "application/json"
const corsHeaderKey string = "Access-Control-Allow-Origin"
const corsHeaderVal string = "*" // unsafe.
const port string = ":8081"

func main() {
	// Set CLI Flags
	setCLIFlags()

	// Setup DB and initiate data fetch
	db.SetupDBFetch()

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
	w.Header().Set(corsHeaderKey, corsHeaderVal)
	var pos models.ISSLocation
	var cor models.Coordinate
	var err error
	pos.Timestamp, cor.Latitude, cor.Longitude, err = db.GetLatestLocation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	pos.Coordinates = cor
	var response []byte
	response, err = json.Marshal(pos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
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
