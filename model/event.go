package model

import (
	"time"
)

type Event struct {
	Model EventModel
}

type EventModel interface {
	GetDateStart() time.Time
}

func (i *Interval) GetDateStart() time.Time { return i.DateStart }

func (l *Lap) GetDateStart() time.Time { return l.DateStart }

func (l *Location) GetDateStart() time.Time { return l.DateStart }

func (p *Position) GetDateStart() time.Time { return p.DateStart }

func (r *RaceControl) GetDateStart() time.Time { return r.DateStart }

func (p *Pit) GetDateStart() time.Time { return p.DateStart }
