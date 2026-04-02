package models

import "time"

type Pit struct {
	Date         time.Time `json:"date"`
	DriverNumber uint      `json:"driver_number"`
	LaneDuration float64   `json:"lane_duration"`
	LapNumber    uint      `json:"lap_number"`
	MeetingKey   uint      `json:"meeting_key"`
	SessionKey   uint      `json:"session_key"`
	StopDuration float64   `json:"stop_duration"`
}
