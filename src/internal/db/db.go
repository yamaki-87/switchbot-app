package db

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

type DBPgxConn struct {
	inner *pgxpool.Pool
	ctx   context.Context
}

var (
	once   sync.Once
	dbConn *DBPgxConn
)

func GetDBConn() *DBPgxConn {
	if dbConn == nil {
		utils.Fatal("Call DbInit pg error")
	}
	return dbConn
}

func DbInit(ctx context.Context, dsn string) error {
	var err error
	once.Do(func() {
		pool, inerr := pgxpool.New(ctx, dsn)
		if inerr != nil {
			err = inerr
		}
		dbConn = &DBPgxConn{
			inner: pool,
			ctx:   ctx,
		}

	})
	return err
}

func (db *DBPgxConn) GetContext() context.Context {
	return db.ctx
}

func (db *DBPgxConn) Query(query string, args ...any) (pgx.Rows, error) {
	return db.inner.Query(db.ctx, query, args...)
}

func (db *DBPgxConn) Exec(query string, args ...any) (pgconn.CommandTag, error) {
	return db.inner.Exec(db.ctx, query, args...)
}

func (db *DBPgxConn) BeginTx() (pgx.Tx, error) {
	return db.inner.Begin(db.ctx)
}

func CoolectRows[T any](rows pgx.Rows) ([]T, error) {
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}
