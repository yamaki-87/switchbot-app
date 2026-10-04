package health

import (
	"context"

	"github.com/yamaki-87/switchbot-app/src/internal/domain"
)

type DeviceListGateway interface {
	GetDeviceList(ctx context.Context) ([]domain.Device, error)
}
