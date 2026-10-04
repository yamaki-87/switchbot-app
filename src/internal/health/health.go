package health

import (
	"context"
	"log/slog"
)

type HealthCheckLogic struct {
	gateway DeviceListGateway
}

func NewHealthCheckLogic(gateway DeviceListGateway) *HealthCheckLogic {
	return &HealthCheckLogic{gateway: gateway}
}

func (h *HealthCheckLogic) Check(ctx context.Context) {
	devices, apiErr := h.gateway.GetDeviceList(ctx)

	if apiErr != nil {
		slog.Error("health check failed", "error", apiErr)
		return
	}

	for _, device := range devices {
		slog.Info("device", "id", device.DeviceID, "name", device.DeviceName, "type", device.DeviceType)
	}
	slog.Info("health check succeeded")
}
