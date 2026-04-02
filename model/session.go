package models

import "time"

type Session struct {
	CircuitKey       uint      `json:"circuit_key"`
	CircuitShortName string    `json:"circuit_short_name"`
	CountryCode      string    `json:"country_code"`
	CountryKey       uint      `json:"country_key"`
	CountryName      string    `json:"country_name"`
	DateEnd          time.Time `json:"date_end"`
	DateStart        time.Time `json:"date_start"`
	GMTOffset        string    `json:"gmt_offset"`
	Location         string    `json:"location"`
	MeetingKey       uint      `json:"meeting_key"`
	SessionKey       uint      `json:"session_key"`
	SessionName      string    `json:"session_name"`
	SessionType      string    `json:"session_type"`
	Year             uint      `json:"year"`
}
