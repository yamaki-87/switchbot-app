package repo

import (
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/dto"
)

type DeviceStatus struct {
	DeviceID         string
	CreateTimeStamp  time.Time
	PowerW           *float64
	VoltageV         *float64
	CurrentMa        *int
	IntervalKwh      *float64
	CollectionStatus int16
}

type DeviceStatusRepo interface {
	InsertBatch(dtos []dto.DeviceStatusInsertDto) error
}

type PostgresDeviceStatusRepo struct {
	db *db.DBPgxConn
}

func NewPostgresDeviceStatusRepo(db *db.DBPgxConn) *PostgresDeviceStatusRepo {
	return &PostgresDeviceStatusRepo{db: db}
}

func (r *PostgresDeviceStatusRepo) InsertBatch(dtos []dto.DeviceStatusInsertDto) error {
	tx, err := r.db.BeginTx()
	if err != nil {
		return err
	}
	ctx := r.db.GetContext()
	defer tx.Rollback(ctx)

	for _, dto := range dtos {
		insertEntity := DeviceStatus{
			DeviceID:         dto.DeviceID,
			CreateTimeStamp:  dto.CreateTimeStamp,
			PowerW:           dto.PowerW,
			VoltageV:         dto.VoltageV,
			CurrentMa:        dto.CurrentMa,
			IntervalKwh:      dto.IntervalKwh,
			CollectionStatus: dto.CollectionStatus,
		}
		_, err := tx.Exec(ctx, `INSERT INTO switchbot_power
			(device_id, create_timestamp, power_w, voltage_v, current_ma, interval_kwh, collection_status)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			insertEntity.DeviceID, insertEntity.CreateTimeStamp, insertEntity.PowerW, insertEntity.VoltageV, insertEntity.CurrentMa, insertEntity.IntervalKwh, insertEntity.CollectionStatus)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
