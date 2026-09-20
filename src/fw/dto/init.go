package dto

import (
	"context"
	"net/http"

	"github.com/yamaki-87/switchbot-app/src/internal/config"
)

type InitOut struct {
	Settings *config.Settings
	Client   *http.Client
	Ctx      context.Context
	Stop     context.CancelFunc
}
