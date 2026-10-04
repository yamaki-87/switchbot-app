package repo

import "time"

type PowerDailySummary struct {
	DeviceID        string    `db:"device_id"`
	ProcessDate     time.Time `db:"process_date"`
	SumKwh          float64   `db:"sum_kwh"`
	CreateTimestamp time.Time `db:"create_timestamp"`
}

type PowerDailySummaryAndNotifyRepo interface {
	InsertBatch(summary []PowerDailySummary, notify []Notify) error
	FindMaxProcessDate() (*time.Time, error)
}
