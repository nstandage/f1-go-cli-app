package model

import (
	"time"

	types "github.com/nstandage/f1-go-cli-app/customtype"
)

type Interval struct {
	DateStart    time.Time             `json:"date"`
	DriverNumber uint                  `json:"driver_number"`
	GapToLeader  *types.FlexibleString `json:"gap_to_leader"`
	Interval     *types.FlexibleString `json:"interval"`
	MeetingKey   uint                  `json:"meeting_key"`
	SessionKey   uint                  `json:"session_key"`
}
