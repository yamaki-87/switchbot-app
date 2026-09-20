package switchbotapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func GetDeviceList(client *http.Client, now time.Time, repDto dto.DeviceListRequest) (*dto.DeviceListBody, error) {
	timestamp := now.UnixMilli()
	nonce := utils.NewNonce()
	sign := utils.CreateSignature(repDto.Token, repDto.Secret, nonce, timestamp)

	url := "https://api.switch-bot.com/v1.1/devices"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	headerReqSet(req, repDto.Token, sign, nonce, timestamp)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deviceListResp dto.DeviceListResponse
	if err := utils.DecodeJSON(resp.Body, &deviceListResp); err != nil {
		return nil, err
	}

	if deviceListResp.StatusCode != SUCCESS_CODE {
		return nil, fmt.Errorf("unexpected status code: %d message: %s", deviceListResp.StatusCode, deviceListResp.Message)
	}

	return &deviceListResp.Body, nil
}
