package switchbotapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func GetDeviceStatus(client *http.Client, now time.Time, reqDto dto.DeviceStatusRequest) (dto.DeviceStatus, error) {
	timestamp := now.UnixMilli()
	nonce := utils.NewNonce()
	sign := utils.CreateSignature(reqDto.Token, reqDto.Secret, nonce, timestamp)

	url := "https://api.switch-bot.com/v1.1/devices/" + reqDto.DeviceID + "/status"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return dto.DeviceStatus{}, err
	}

	headerReqSet(req, reqDto.Token, sign, nonce, timestamp)

	resp, err := client.Do(req)
	if err != nil {
		return dto.DeviceStatus{}, err
	}
	defer resp.Body.Close()

	var result dto.DeviceStatusResponse

	if err := utils.DecodeJSON(resp.Body, &result); err != nil {
		return dto.DeviceStatus{}, err
	}

	if result.StatusCode != SUCCESS_CODE {
		return dto.DeviceStatus{}, fmt.Errorf("unexpected status code: %d message: %s", result.StatusCode, result.Message)
	}

	return result.Body, nil
}
