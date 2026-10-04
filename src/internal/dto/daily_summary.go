package dto

import "time"

type SelectDailySummaryIn struct {
	StartDate time.Time
	EndDate   time.Time
}
