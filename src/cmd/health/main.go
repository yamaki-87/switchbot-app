package main

import (
	"log/slog"

	"github.com/yamaki-87/switchbot-app/src/fw"
	"github.com/yamaki-87/switchbot-app/src/internal/health"
	"github.com/yamaki-87/switchbot-app/src/internal/switchbotapi"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func main() {
	initOut, err := fw.Init()
	utils.FatalIfErr(err, "init failed")
	defer initOut.Stop()
	gateway := switchbotapi.NewClient(initOut.Client, initOut.Settings.Secret.Token, initOut.Settings.Secret.Secret)

	slog.Info("helath check start")
	healthCheck := health.NewHealthCheckLogic(gateway)
	healthCheck.Check(initOut.Ctx)
	slog.Info("helath check end")
}
