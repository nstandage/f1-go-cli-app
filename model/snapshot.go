package model

import (
	"time"
)

type Snapshot struct {
	SessionBar      *SessionBarSnapShot
	RaceControlMsgs []string
	Drivers         []DriverSnapshot
	LastLap         []string
	LastLapIsPitOut []bool
	Intervals       []string
	GapsToLeaders   []string
	PitCounts       []string
	TireCompounds   []string
	TireAges        []string
}

type SessionBarSnapShot struct {
	EventName        string
	EventType        string
	FastestLap		 *FastestLapSnapshot
	CurrentLap       uint
	TotalLaps        uint
	IsReplay         bool
	EventDate        time.Time
}

type FastestLapSnapshot struct {
	LapTime   string
	Driver string
	LapNumber string
}


type DriverSnapshot struct {
	Name string
	IsFastestLap bool
}
