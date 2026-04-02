package models

type Location struct {
	Date         string `json:"date"`
	DriverNumber uint   `json:"driver_number"`
	MeetingKey   uint   `json:"meeting_key"`
	SessionKey   uint   `json:"session_key"`
	X            uint   `json:"x"`
	Y            uint   `json:"y"`
	Z            uint   `json:"z"`
}
