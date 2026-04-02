package models

type SessionData struct {
	Meeting   Meeting
	Session   Session
	Drivers   []Driver
	Intervals []Interval
	Laps      []Lap
	Locations []Location
	Pits      []Pit
	Positions []Position
	Stints    []Stint
}
