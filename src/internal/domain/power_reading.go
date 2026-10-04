package domain

type PowerReading struct {
	DeviceID string
	PowerW float64
	VoltageV float64
	CurrentMa int
}

type Device struct {
	DeviceID string
	DeviceName string
	DeviceType string
}
