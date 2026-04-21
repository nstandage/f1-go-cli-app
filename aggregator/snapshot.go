package aggregator

import (
	"time"
)

type Snapshot struct {
	SessionBar      SessionBarSnapShot
	RaceControlMsgs []string
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
