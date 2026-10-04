package postgres

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
)

var _ repo.PowerDailySummaryAndNotifyRepo = (*PostgresPowerDailySummaryAndNotifyRepo)(nil)

type PostgresPowerDailySummaryAndNotifyRepo struct {
	db *db.DBPgxConn
}

func NewPostgresPowerDailySummaryRepo(db *db.DBPgxConn) *PostgresPowerDailySummaryAndNotifyRepo {
	return &PostgresPowerDailySummaryAndNotifyRepo{db: db}
}

func (r *PostgresPowerDailySummaryAndNotifyRepo) InsertBatch(summary []repo.PowerDailySummary, notify []repo.Notify) error {
	tx, err := r.db.BeginTx()
	if err != nil {
		return err
	}

	ctx := r.db.GetContext()
	defer tx.Rollback(ctx)

	err = r.insertSwitchbotPowerSummary(summary, tx, ctx)
	if err != nil {
		return err
	}

	err = r.insertSwitchbotNotify(notify, tx, ctx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresPowerDailySummaryAndNotifyRepo) insertSwitchbotPowerSummary(summary []repo.PowerDailySummary, tx pgx.Tx, ctx context.Context) error {
	for _, s := range summary {
		_, err := tx.Exec(ctx, `INSERT INTO switchbot_power_daily_summary
			(device_id, process_date, sum_kwh) VALUES ($1, $2, $3)`,
			s.DeviceID, s.ProcessDate, s.SumKwh)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresPowerDailySummaryAndNotifyRepo) insertSwitchbotNotify(notify []repo.Notify, tx pgx.Tx, ctx context.Context) error {
	for _, n := range notify {
		_, err := tx.Exec(ctx, `INSERT INTO switchbot_notify
			(device_id, status, notify_msg) VALUES ($1, $2, $3)`,
			n.DeviceId, n.Status, n.NotifyMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresPowerDailySummaryAndNotifyRepo) FindMaxProcessDate() (*time.Time, error) {
	var maxDate *time.Time
	err := r.db.QueryRow(`SELECT MAX(process_date) FROM switchbot_power_daily_summary`).Scan(&maxDate)
	if err != nil {
		return nil, err
	}
	return maxDate, nil
}
