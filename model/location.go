package model

import "time"

type Location struct {
	DateStart    time.Time `json:"date"`
	DriverNumber uint      `json:"driver_number"`
	MeetingKey   uint      `json:"meeting_key"`
	SessionKey   uint      `json:"session_key"`
	X            int       `json:"x"`
	Y            int       `json:"y"`
	Z            int       `json:"z"`
}
