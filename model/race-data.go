package model

import "time"

type RaceData struct {
	Meeting      *Meeting
	Session      *Session
	Drivers      []Driver
	Stints       []Stint
	TotalLaps    uint
	StartingGrid []StartingGrid
	SessionStart time.Time
}
