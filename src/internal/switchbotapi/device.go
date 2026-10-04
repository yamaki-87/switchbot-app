package switchbotapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/domain"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func (c *Client) GetDeviceStatus(ctx context.Context, deviceID string) (domain.PowerReading, error) {
	timestamp := time.Now().UnixMilli()
	nonce := utils.NewNonce()
	sign := utils.CreateSignature(c.token, c.secret, nonce, timestamp)

	url := "https://api.switch-bot.com/v1.1/devices/" + deviceID + "/status"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.PowerReading{}, err
	}

	headerReqSet(req, c.token, sign, nonce, timestamp)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.PowerReading{}, err
	}
	defer resp.Body.Close()

	var result DeviceStatusResponse

	if err := utils.DecodeJSON(resp.Body, &result); err != nil {
		return domain.PowerReading{}, err
	}

	if result.StatusCode != SUCCESS_CODE {
		return domain.PowerReading{}, fmt.Errorf("unexpected status code: %d message: %s", result.StatusCode, result.Message)
	}

	return domain.PowerReading{
		DeviceID: result.Body.DeviceID,
		PowerW: result.Body.Weight,
		VoltageV: result.Body.Voltage,
		CurrentMa: result.Body.ElectricCurrent,
	}, nil
}
