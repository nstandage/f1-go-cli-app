package models

type Position struct {
	Date         string `json:"date"`
	DriverNumber uint   `json:"driver_number"`
	MeetingKey   uint   `json:"meeting_key"`
	Position     uint   `json:"position"`
	SessionKey   uint   `json:"session_key"`
}
