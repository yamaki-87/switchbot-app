package collector

import (
	"context"
	"sync"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/domain"
	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

type CollectorLogic struct {
	deviceStatusRepo repo.DeviceStatusRepo
	deviceStatusGateway DeviceStatusGateway
	previousHolder   map[string]dto.PreviousCollector
	mu               sync.Mutex
}

func NewCollectorLogic(deviceStatusRepo repo.DeviceStatusRepo, gateway DeviceStatusGateway) *CollectorLogic {
	return &CollectorLogic{
		deviceStatusRepo: deviceStatusRepo,
		deviceStatusGateway: gateway,
		previousHolder:   make(map[string]dto.PreviousCollector),
		mu:               sync.Mutex{},
	}
}

func (c *CollectorLogic) Collect(ctx context.Context, deviceID string) (dto.DeviceStatusInsertDto, error) {
	now := utils.GetTimeNow()
	deviceStatus, apiErr := c.deviceStatusGateway.GetDeviceStatus(ctx, deviceID)

	var intervalKwh float64

	if apiErr == nil {
		intervalKwh = c.updatePrevious(deviceID, now, deviceStatus)
	}

	collectionStatus := getCollectionStatus(apiErr)

	insertDto := deviceStatusInsertDtoTo(
		&deviceStatus,
		now,
		intervalKwh,
		collectionStatus,
		deviceID,
	)

	return insertDto, apiErr
}

func (c *CollectorLogic) InsertBatch(in []dto.DeviceStatusInsertDto) error {
	err := c.deviceStatusRepo.InsertBatch(in)
	if err != nil {
		return err
	}
	return nil
}

func getCollectionStatus(apiErr error) CollectionStatus {
	if apiErr != nil {
		return FAILED
	}
	return SUCCESS
}

func deviceStatusInsertDtoTo(deviceStatus *domain.PowerReading, now time.Time, intervalKwh float64, collectionStatus CollectionStatus, deviceId string) dto.DeviceStatusInsertDto {
	var result dto.DeviceStatusInsertDto
	if collectionStatus == FAILED {
		result = dto.DeviceStatusInsertDto{
			// API取得時にDeviceStatusが取得できないため、予め取得したdeviceIdを使用
			DeviceID:         deviceId,
			CreateTimeStamp:  now,
			PowerW:           nil,
			VoltageV:         nil,
			CurrentMa:        nil,
			IntervalKwh:      nil,
			CollectionStatus: int16(collectionStatus),
		}
	} else {
		result = dto.DeviceStatusInsertDto{
			DeviceID:         deviceStatus.DeviceID,
			CreateTimeStamp:  now,
			PowerW:           &deviceStatus.PowerW,
			VoltageV:         &deviceStatus.VoltageV,
			CurrentMa:        &deviceStatus.CurrentMa,
			IntervalKwh:      &intervalKwh,
			CollectionStatus: int16(collectionStatus),
		}
	}

	return result
}

func (c *CollectorLogic) updatePrevious(
	deviceID string,
	now time.Time,
	deviceStatus domain.PowerReading,
) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	var intervalKwh float64

	if previousValue, ok := c.previousHolder[deviceID]; ok {
		elapsed := now.Sub(previousValue.PreviousTime)
		averageW := (previousValue.PreviousStatus.PowerW + deviceStatus.PowerW) / 2
		intervalKwh = averageW * elapsed.Hours() / 1000
	}

	c.previousHolder[deviceID] = dto.PreviousCollector{
		PreviousTime:   now,
		PreviousStatus: &deviceStatus,
	}

	return intervalKwh
}
