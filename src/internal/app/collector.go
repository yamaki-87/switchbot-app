package app

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/collector"
	"github.com/yamaki-87/switchbot-app/src/internal/dto"
)

type CollectorApp struct {
	collectorLogic *collector.CollectorLogic
	devices        []dto.DeviceMaster
}

func NewCollectorApp(collectorLogic *collector.CollectorLogic, devices []dto.DeviceMaster) *CollectorApp {
	return &CollectorApp{
		collectorLogic: collectorLogic,
		devices:        devices,
	}
}

func (a *CollectorApp) Run(ctx context.Context, interval time.Duration) {
	a.outputStartLog(interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.collectOnce(ctx)

		case <-ctx.Done():
			slog.Info("context done")
			return
		}
	}
}

func (a *CollectorApp) outputStartLog(interval time.Duration) {
	deviceNames := make([]string, 0, len(a.devices))
	for _, device := range a.devices {
		deviceNames = append(deviceNames, device.DeviceName)
	}
	result := strings.Join(deviceNames, ", ")
	slog.Info("collector started", " devices", result, "interval", interval)
}

func (a *CollectorApp) collectOnce(ctx context.Context) {
	resultCh := make(
		chan dto.DeviceStatusInsertDto,
		len(a.devices),
	)

	for _, device := range a.devices {
		go func() {
			status, err := a.collectorLogic.Collect(ctx, device.DeviceId)
			if err != nil {
				slog.Error(
					"collect failed",
					"device_id", device.DeviceId,
					"err", err,
				)
			}

			resultCh <- status
		}()
	}

	results := make(
		[]dto.DeviceStatusInsertDto,
		0,
		len(a.devices),
	)

	for range a.devices {
		results = append(results, <-resultCh)
	}

	if err := a.collectorLogic.InsertBatch(results); err != nil {
		slog.Error("insert failed", "err", err)
	}
}
