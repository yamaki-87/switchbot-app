package repo

import "github.com/yamaki-87/switchbot-app/src/internal/db"

type DeviceMaster struct {
	DeviceId           string `db:"device_id"`
	DeviceName         string `db:"device_name"`
	PollingIntervalSec int16  `db:"polling_interval_sec"`
}

type DeviceMasterRepo interface {
	FindDevicesNotDeleted() ([]DeviceMaster, error)
}

type PostgresDeviceMasterRepo struct {
	db *db.DBPgxConn
}

func NewPostgresDeviceMasterRepo(db *db.DBPgxConn) *PostgresDeviceMasterRepo {
	return &PostgresDeviceMasterRepo{db: db}
}

func (r *PostgresDeviceMasterRepo) FindDevicesNotDeleted() ([]DeviceMaster, error) {
	rows, err := r.db.Query(`
SELECT
    SD.DEVICE_ID
    , SD.DEVICE_NAME
    , SDS.POLLING_INTERVAL_SEC
        FROM
    SWITCHBOT_DEVICE AS SD
    INNER JOIN SWITCHBOT_DEVICE_SETTINGS AS SDS
        ON SD.DEVICE_ID = SDS.DEVICE_ID
WHERE
    SD.IS_DELETED = $1
`, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices, err := db.CoolectRows[DeviceMaster](rows)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
