package devicemaster

import (
	"fmt"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
)

type DeviceMasterManager struct {
	deviceMaster     []dto.DeviceMaster
	deviceMasterRepo repo.DeviceMasterRepo
}

func NewDeviceMasterManager(deviceMasterRepo repo.DeviceMasterRepo) (*DeviceMasterManager, error) {
	deviceMaster, err := deviceMasterRepo.FindDevicesNotDeleted()
	if err != nil {
		return nil, err
	}

	if len(deviceMaster) == 0 {
		return nil, fmt.Errorf("find device not deleted failed")
	}

	return &DeviceMasterManager{
		deviceMaster:     deviceMasterToDtos(deviceMaster),
		deviceMasterRepo: deviceMasterRepo,
	}, nil
}

func (d *DeviceMasterManager) GetDeviceMaster() []dto.DeviceMaster {
	return d.deviceMaster
}

func deviceMasterToDtos(deviceMaster []repo.DeviceMaster) []dto.DeviceMaster {
	dtoList := make([]dto.DeviceMaster, 0, len(deviceMaster))
	for _, dm := range deviceMaster {
		dtoList = append(dtoList, deviceMasterToDto(dm))
	}
	return dtoList
}

func deviceMasterToDto(deviceMaster repo.DeviceMaster) dto.DeviceMaster {
	return dto.DeviceMaster{
		DeviceId:           deviceMaster.DeviceId,
		DeviceName:         deviceMaster.DeviceName,
		PollingIntervalSec: deviceMaster.PollingIntervalSec,
	}
}
