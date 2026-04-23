package model

import (
	"time"
)

type Snapshot struct {
	SessionBar      *SessionBarSnapShot
	RaceControlMsgs []string
	DriverNames     []string
	LastLap         []string
}

type SessionBarSnapShot struct {
	EventName        string
	EventType        string
	CurrentLap       uint
	TotalLaps        uint
	FastestLapTime   time.Time
	FastestLapDriver string
	FastestLapNumber uint
	IsReplay         bool
	EventDate        time.Time
}
