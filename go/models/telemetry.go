package models

import "gorm.io/gorm"

type Telemetry struct {
	TripId    string  `json:"trip_id"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Speed     float64 `json:"speed"`
	Timestamp string  `json:"timestamp"`
}

func MigrateTelemetry(db *gorm.DB) error {
	err := db.AutoMigrate(&Telemetry{})
	return err
}