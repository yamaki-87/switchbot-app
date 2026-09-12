package collector

import (
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
	"github.com/yamaki-87/switchbot-app/src/internal/switchbotapi"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

type CollectorLogic struct {
	deviceStatusRepo repo.DeviceStatusRepo
}

func NewCollectorLogic(deviceStatusRepo repo.DeviceStatusRepo) *CollectorLogic {
	return &CollectorLogic{
		deviceStatusRepo: deviceStatusRepo,
	}
}

var previousManager = make(map[string]dto.PreviousCollector)

func (c *CollectorLogic) Collect(in *dto.CollectorIn) (dto.DeviceStatusInsertDto, error) {
	now := utils.GetTimeNow()
	req := dto.DeviceStatusRequest{
		Token:    in.Token,
		Secret:   in.Secret,
		DeviceID: in.DeviceId,
	}
	deviceStatus, apiErr := switchbotapi.GetDeviceStatus(in.Client, now, req)

	previousValue, ok := previousManager[in.DeviceId]
	var intervalKwh float64
	if ok {
		elapsed := now.Sub(previousValue.PreviousTime)
		averageW := (previousValue.PreviousStatus.Weight + deviceStatus.Weight) / 2
		intervalKwh = averageW * elapsed.Hours() / 1000
	}
	collectionStatus := getCollectionStatus(apiErr)
	insertDto := deviceStatusInsertDtoTo(&deviceStatus, now, intervalKwh, collectionStatus)

	previousManager[in.DeviceId] = dto.PreviousCollector{
		PreviousTime:   now,
		PreviousStatus: &deviceStatus,
	}
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

func deviceStatusInsertDtoTo(deviceStatus *dto.DeviceStatus, now time.Time, intervalKwh float64, collectionStatus CollectionStatus) dto.DeviceStatusInsertDto {
	var result dto.DeviceStatusInsertDto
	if collectionStatus == FAILED {
		result = dto.DeviceStatusInsertDto{
			DeviceID:         deviceStatus.DeviceID,
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
