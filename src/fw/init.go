package fw

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/yamaki-87/switchbot-app/src/fw/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/config"
)

func Init() (*dto.InitOut, error) {
	_ = godotenv.Load()

	err := config.LoadSettings(os.Getenv("SETTING_PATH"))
	if err != nil {
		return nil, err
	}

	settings := config.GetSettings()
	client := &http.Client{
		Timeout: time.Duration(settings.Config.App.IntervalTimeout) * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	return &dto.InitOut{
		Settings: settings,
		Client:   client,
		Ctx:      ctx,
		Stop:     stop,
	}, nil
}
