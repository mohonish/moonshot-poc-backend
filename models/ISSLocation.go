package models

import (
	"time"
)

type ISSLocation struct {
	Timestamp	time.Time	`json:"timestamp"`
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"longitude"`
}