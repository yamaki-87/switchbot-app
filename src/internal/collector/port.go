package collector

import (
	"context"

	"github.com/yamaki-87/switchbot-app/src/internal/domain"
)

type DeviceStatusGateway interface {
	GetDeviceStatus(ctx context.Context, deviceID string) (domain.PowerReading, error)
}
