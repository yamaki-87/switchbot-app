package dto

import "time"

type DeviceStatusInsertDto struct {
	DeviceID         string
	CreateTimeStamp  time.Time
	PowerW           *float64
	VoltageV         *float64
	CurrentMa        *int
	IntervalKwh      *float64
	CollectionStatus int16
}
