package repo

import (
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/status"
)

type Notify struct {
	NotifyId        int                 `db:"notify_id"`
	DeviceId        string              `db:"device_id"`
	Status          status.NotifyStatus `db:"status"`
	NotifyMsg       string              `db:"notify_msg"`
	CreateTimeStamp time.Time           `db:"create_timestamp"`
	UpdateTimeStamp time.Time           `db:"update_timestamp"`
}
