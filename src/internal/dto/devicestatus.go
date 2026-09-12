package dto

import "time"

type DeviceStatus struct {
	Power            string  `json:"power"`
	Voltage          float64 `json:"voltage"`
	Weight           float64 `json:"weight"`
	ElectricityOfDay int     `json:"electricityOfDay"`
	ElectricCurrent  int     `json:"electricCurrent"`
	DeviceID         string  `json:"deviceId"`
	DeviceType       string  `json:"deviceType"`
}

type DeviceStatusRequest struct {
	Token    string `json:"token"`
	Secret   string `json:"secret"`
	DeviceID string `json:"deviceId"`
}

type DeviceStatusResponse struct {
	StatusCode int          `json:"statusCode"`
	Message    string       `json:"message"`
	Body       DeviceStatus `json:"body"`
}

type DeviceStatusInsertDto struct {
	DeviceID         string
	CreateTimeStamp  time.Time
	PowerW           *float64
	VoltageV         *float64
	CurrentMa        *int
	IntervalKwh      *float64
	CollectionStatus int16
}
