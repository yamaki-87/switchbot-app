package dto

import "net/http"

type DeviceListResponse struct {
	StatusCode int            `json:"statusCode"`
	Body       DeviceListBody `json:"body"`
	Message    string         `json:"message"`
}

type DeviceListBody struct {
	DeviceList         []Device         `json:"deviceList"`
	InfraredRemoteList []InfraredRemote `json:"infraredRemoteList"`
}

type Device struct {
	DeviceID           string `json:"deviceId"`
	DeviceName         string `json:"deviceName"`
	DeviceType         string `json:"deviceType"`
	EnableCloudService bool   `json:"enableCloudService"`
	HubDeviceID        string `json:"hubDeviceId"`
}

type InfraredRemote struct {
	DeviceID    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	RemoteType  string `json:"remoteType"`
	HubDeviceID string `json:"hubDeviceId"`
}

type DeviceListRequest struct {
	Token  string
	Secret string
}

type HealthCheckIn struct {
	Client *http.Client
	Token  string
	Secret string
}
