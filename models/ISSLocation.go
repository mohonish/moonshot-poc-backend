package models

import (
	"time"
)

// Example response:
// {"iss_position": {"longitude": "48.4187", "latitude": "-43.3559"}, "message": "success", "timestamp": 1590903925}

type Coordinate struct {
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"longitude"`
}

type ISSLocation struct {
	Coordinates	Coordinate	`json:"iss_position"`
	Timestamp	time.Time	`json:"timestamp"`
}