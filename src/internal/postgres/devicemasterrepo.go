package postgres

import (
	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
)

var _ repo.DeviceMasterRepo = (*PostgresDeviceMasterRepo)(nil)

type PostgresDeviceMasterRepo struct {
	db *db.DBPgxConn
}

func NewPostgresDeviceMasterRepo(db *db.DBPgxConn) *PostgresDeviceMasterRepo {
	return &PostgresDeviceMasterRepo{db: db}
}

func (r *PostgresDeviceMasterRepo) FindDevicesNotDeleted() ([]repo.DeviceMaster, error) {
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

	devices, err := db.CoolectRows[repo.DeviceMaster](rows)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
