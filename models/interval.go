package models

type Interval struct {
	Date         string  `json:"date"`
	DriverNumber uint    `json:"driver_number"`
	GapToLeader  float64 `json:"gap_to_leader"`
	Interval     float64 `json:"interval"`
	MeetingKey   uint    `json:"meeting_key"`
	SessionKey   uint    `json:"session_key"`
}
