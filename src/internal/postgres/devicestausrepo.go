package postgres

import (
	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
)

var _ repo.DeviceStatusRepo = (*PostgresDeviceStatusRepo)(nil)

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
		insertEntity := repo.DeviceStatus{
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

func (r *PostgresDeviceStatusRepo) SelectDailySummary(dto dto.SelectDailySummaryIn) ([]repo.SelectDailySummary, error) {
	tx, err := r.db.BeginTx()
	if err != nil {
		return nil, err
	}
	ctx := r.db.GetContext()
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT device_id, SUM(interval_kwh) AS sum_kwh
		FROM switchbot_power
		WHERE CREATE_TIMESTAMP >= $1 AND CREATE_TIMESTAMP < $2
		GROUP BY device_id`, dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := db.CoolectRows[repo.SelectDailySummary](rows)
	if err != nil {
		return nil, err
	}
	return results, nil
}
