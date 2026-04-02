package model

type EventType int

const (
	DriverEventType EventType = iota
	IntervalEventType
	LapEventType
	LocationEventType
	MeetingEventType
	PitEventType
	PositionEventType
	RaceControlEventType
	SessionEventType
	StintEventType
)

type Event interface {
	EventType() EventType
}

type DriverEvent struct {
	data Driver
}

func (d *DriverEvent) EventType() EventType { return DriverEventType }

type IntervalEvent struct {
	data Interval
}

func (i *IntervalEvent) EventType() EventType { return IntervalEventType }

type LapEvent struct {
	data Lap
}

func (l *LapEvent) EventType() EventType { return LapEventType }

type LocationEvent struct {
	data Location
}

func (l *LocationEvent) EventType() EventType { return LocationEventType }

type MeetingEvent struct {
	data Meeting
}

func (m *MeetingEvent) EventType() EventType { return MeetingEventType }

type PitEvent struct {
	data Pit
}

func (p *PitEvent) EventType() EventType { return PitEventType }

type PositionEvent struct {
	data Position
}

func (p *PositionEvent) EventType() EventType { return PositionEventType }

type RaceControlEvent struct {
	data RaceControl
}

func (r *RaceControlEvent) EventType() EventType { return RaceControlEventType }

type SessionEvent struct {
	data Session
}

func (s *SessionEvent) EventType() EventType { return SessionEventType }

type StintEvent struct {
	data Stint
}

func (s *StintEvent) EventType() EventType { return StintEventType }
