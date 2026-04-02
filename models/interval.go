package models

import (
	"time"

	"github.com/nstandage/f1-go-cli-app/types"
)

type Interval struct {
	Date         time.Time             `json:"date"`
	DriverNumber uint                  `json:"driver_number"`
	GapToLeader  *types.FlexibleString `json:"gap_to_leader"`
	Interval     *types.FlexibleString `json:"interval"`
	MeetingKey   uint                  `json:"meeting_key"`
	SessionKey   uint                  `json:"session_key"`
}
