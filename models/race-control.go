package models

import "time"

type RaceControl struct {
	Category        string    `json:"category"`
	Date            time.Time `json:"date"`
	DriverNumber    int       `json:"driver_number"`
	Flag            string    `json:"flag"`
	LapNumber       uint      `json:"lap_number"`
	MeetingKey      uint      `json:"meeting_key"`
	Message         string    `json:"message"`
	QualifyingPhase *uint     `json:"qualifying_phase"`
	Scope           string    `json:"scope"`
	Sector          *uint     `json:"sector"`
	SessionKey      int       `json:"session_key"`
}
