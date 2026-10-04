package dto

import (
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/domain"
)

type PreviousCollector struct {
	PreviousTime time.Time
	PreviousStatus *domain.PowerReading
}
