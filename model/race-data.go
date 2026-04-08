package model

type RaceData struct {
	Meeting *Meeting
	Session *Session
	Drivers []Driver
	Stints  []Stint
}
