package repo

type DeviceMaster struct {
	DeviceId           string `db:"device_id"`
	DeviceName         string `db:"device_name"`
	PollingIntervalSec int16  `db:"polling_interval_sec"`
}

type DeviceMasterRepo interface {
	FindDevicesNotDeleted() ([]DeviceMaster, error)
}
