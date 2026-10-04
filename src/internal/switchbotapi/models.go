package switchbotapi

type DeviceStatus struct {
	Power            string  `json:"power"`
	Voltage          float64 `json:"voltage"`
	Weight           float64 `json:"weight"`
	ElectricityOfDay int     `json:"electricityOfDay"`
	ElectricCurrent  int     `json:"electricCurrent"`
	DeviceID         string  `json:"deviceId"`
	DeviceType       string  `json:"deviceType"`
}

type DeviceStatusResponse struct {
	StatusCode int          `json:"statusCode"`
	Message    string       `json:"message"`
	Body       DeviceStatus `json:"body"`
}

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
