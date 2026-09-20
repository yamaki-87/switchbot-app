package health

import (
	"log/slog"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/switchbotapi"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

type HealthCheckLogic struct {
}

func NewHealthCheckLogic() *HealthCheckLogic {
	return &HealthCheckLogic{}
}

func (h *HealthCheckLogic) Check(in *dto.HealthCheckIn) {
	now := utils.GetTimeNow()
	req := dto.DeviceListRequest{
		Token:  in.Token,
		Secret: in.Secret,
	}
	devices, apiErr := switchbotapi.GetDeviceList(in.Client, now, req)
	if apiErr != nil {
		slog.Error("health check failed", "error", apiErr)
		return
	}

	for _, device := range devices.DeviceList {
		slog.Info("device", "id", device.DeviceID, "name", device.DeviceName, "type", device.DeviceType)
	}
	slog.Info("health check succeeded")
}
