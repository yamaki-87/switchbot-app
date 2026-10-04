package repo

import (
	"time"
	"github.com/yamaki-87/switchbot-app/src/internal/dto"
)

type DeviceStatus struct {
	DeviceID         string
	CreateTimeStamp  time.Time
	PowerW           *float64
	VoltageV         *float64
	CurrentMa        *int
	IntervalKwh      *float64
	CollectionStatus int16
}

type SelectDailySummary struct {
	DeviceID string  `db:"device_id"`
	SumKwh   float64 `db:"sum_kwh"`
}

type DeviceStatusRepo interface {
	InsertBatch(dtos []dto.DeviceStatusInsertDto) error
	SelectDailySummary(dto dto.SelectDailySummaryIn) ([]SelectDailySummary, error)
}
