package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mohonish/moonshot-backend/version"
	"log"
	"net/http"
)

const contentTypeKey string = "Content-Type"
const contentTypeVal string = "application/json"

func main() {
	// Set CLI Flags
	setCLIFlags()

	// Create mux base router
	r := mux.NewRouter()

	// Create mux subrouter
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/iss/coordinates", baseRouteGET).Methods(http.MethodGet)
	api.HandleFunc("", routeNotFound)

	// Listen on 8081
	log.Fatal(http.ListenAndServe(":8081", r))
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
