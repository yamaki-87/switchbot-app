package collector

import (
	"sync"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
	"github.com/yamaki-87/switchbot-app/src/internal/switchbotapi"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

type CollectorLogic struct {
	deviceStatusRepo repo.DeviceStatusRepo
	previousHolder   map[string]dto.PreviousCollector
	mu               sync.Mutex
}

func NewCollectorLogic(deviceStatusRepo repo.DeviceStatusRepo) *CollectorLogic {
	return &CollectorLogic{
		deviceStatusRepo: deviceStatusRepo,
		previousHolder:   make(map[string]dto.PreviousCollector),
		mu:               sync.Mutex{},
	}
}

func (c *CollectorLogic) Collect(in *dto.CollectorIn) (dto.DeviceStatusInsertDto, error) {
	now := utils.GetTimeNow()

	req := dto.DeviceStatusRequest{
		Token:    in.Token,
		Secret:   in.Secret,
		DeviceID: in.DeviceId,
	}

	deviceStatus, apiErr := switchbotapi.GetDeviceStatus(
		in.Client,
		now,
		req,
	)

	var intervalKwh float64

	if apiErr == nil {
		intervalKwh = c.updatePrevious(in.DeviceId, now, deviceStatus)
	}

	collectionStatus := getCollectionStatus(apiErr)

	insertDto := deviceStatusInsertDtoTo(
		&deviceStatus,
		now,
		intervalKwh,
		collectionStatus,
		in.DeviceId,
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

func deviceStatusInsertDtoTo(deviceStatus *dto.DeviceStatus, now time.Time, intervalKwh float64, collectionStatus CollectionStatus, deviceId string) dto.DeviceStatusInsertDto {
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
			PowerW:           &deviceStatus.Weight,
			VoltageV:         &deviceStatus.Voltage,
			CurrentMa:        &deviceStatus.ElectricCurrent,
			IntervalKwh:      &intervalKwh,
			CollectionStatus: int16(collectionStatus),
		}
	}

	return result
}

func (c *CollectorLogic) updatePrevious(
	deviceID string,
	now time.Time,
	deviceStatus dto.DeviceStatus,
) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	var intervalKwh float64

	if previousValue, ok := c.previousHolder[deviceID]; ok {
		elapsed := now.Sub(previousValue.PreviousTime)
		averageW := (previousValue.PreviousStatus.Weight + deviceStatus.Weight) / 2
		intervalKwh = averageW * elapsed.Hours() / 1000
	}

	c.previousHolder[deviceID] = dto.PreviousCollector{
		PreviousTime:   now,
		PreviousStatus: &deviceStatus,
	}

	return intervalKwh
}
