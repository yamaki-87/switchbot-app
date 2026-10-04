package switchbotapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/domain"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func (c *Client) GetDeviceList(ctx context.Context) ([]domain.Device, error) {
	timestamp := time.Now().UnixMilli()
	nonce := utils.NewNonce()
	sign := utils.CreateSignature(c.token, c.secret, nonce, timestamp)

	url := "https://api.switch-bot.com/v1.1/devices"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	headerReqSet(req, c.token, sign, nonce, timestamp)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deviceListResp DeviceListResponse
	if err := utils.DecodeJSON(resp.Body, &deviceListResp); err != nil {
		return nil, err
	}

	if deviceListResp.StatusCode != SUCCESS_CODE {
		return nil, fmt.Errorf("unexpected status code: %d message: %s", deviceListResp.StatusCode, deviceListResp.Message)
	}

	devices := make([]domain.Device, 0, len(deviceListResp.Body.DeviceList))
	for _, device := range deviceListResp.Body.DeviceList {
		devices = append(devices, domain.Device{DeviceID: device.DeviceID, DeviceName: device.DeviceName, DeviceType: device.DeviceType})
	}
	return devices, nil
}
