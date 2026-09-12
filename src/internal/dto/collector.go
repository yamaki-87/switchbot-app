package dto

import (
	"net/http"
	"time"
)

type CollectorIn struct {
	Client   *http.Client
	DeviceId string
	Token    string
	Secret   string
}

type PreviousCollector struct {
	PreviousTime   time.Time
	PreviousStatus *DeviceStatus
}
