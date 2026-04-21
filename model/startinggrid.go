package model

type StartingGrid struct {
	Position     uint    `json:"position"`
	DriverNumber uint    `json:"driver_number"`
	LapDuration  float64 `json:"lap_duration"`
	MeetingKey   uint    `json:"meeting_key"`
	SessionKey   uint    `json:"session_key"`
}
