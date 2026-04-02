package models

type Lap struct {
	DateStart       string  `json:"date_start"`
	DriverNumber    uint    `json:"driver_number"`
	DurationSector1 float64 `json:"duration_sector_1"`
	DurationSector2 float64 `json:"duration_sector_2"`
	DurationSector3 float64 `json:"duration_sector_3"`
	I1Speed         uint    `json:"i1_speed"`
	I2Speed         uint    `json:"i2_speed"`
	IsPitOutLap     bool    `json:"is_pit_out_lap"`
	LapDuration     float64 `json:"lap_duration"`
	LapNumber       uint    `json:"lap_number"`
	MeetingKey      uint    `json:"meeting_key"`
	SegmentsSector1 []uint  `json:"segments_sector_1"`
	SegmentsSector2 []uint  `json:"segments_sector_2"`
	SegmentsSector3 []uint  `json:"segments_sector_3"`
	SessionKey      uint    `json:"session_key"`
	StSpeed         uint    `json:"st_speed"`
}
