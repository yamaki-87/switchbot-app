package main

import (
	"log/slog"

	"github.com/yamaki-87/switchbot-app/src/fw"
	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/health"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func main() {
	initOut, err := fw.Init()
	utils.FatalIfErr(err, "init failed")
	defer initOut.Stop()

	slog.Info("helath check start")
	healthCheck := health.NewHealthCheckLogic()
	healthCheck.Check(&dto.HealthCheckIn{
		Client: initOut.Client,
		Token:  initOut.Settings.Secret.Token,
		Secret: initOut.Settings.Secret.Secret,
	})
	slog.Info("helath check end")
}
