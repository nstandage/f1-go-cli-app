package model

type Stint struct {
	Compound       string `json:"compound"`
	DriverNumber   uint   `json:"driver_number"`
	LapEnd         uint   `json:"lap_end"`
	LapStart       uint   `json:"lap_start"`
	MeetingKey     uint   `json:"meeting_key"`
	SessionKey     uint   `json:"session_key"`
	StintNumber    uint   `json:"stint_number"`
	TyreAgeAtStart uint   `json:"tyre_age_at_start"`
}
