package models

type Driver struct {
	BroadcastName string `json:"broadcast_name"`
	DriverNumber  uint   `json:"driver_number"`
	FirstName     string `json:"first_name"`
	FullName      string `json:"full_name"`
	HeadshotURL   string `json:"headshot_url"`
	LastName      string `json:"last_name"`
	MeetingKey    uint   `json:"meeting_key"`
	NameAcronym   string `json:"name_acronym"`
	SessionKey    uint   `json:"session_key"`
	TeamColour    string `json:"team_colour"`
	TeamName      string `json:"team_name"`
}
