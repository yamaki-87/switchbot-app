package main

import (
	"log/slog"

	"github.com/yamaki-87/switchbot-app/src/fw"
	"github.com/yamaki-87/switchbot-app/src/internal/dailyas"
	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/postgres"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func main() {
	initOut, err := fw.Init()
	utils.FatalIfErr(err, "init failed")
	defer initOut.Stop()

	err = db.DbInit(initOut.Ctx, initOut.Settings.Config.Database.URL)
	utils.FatalIfErr(err, "db init failed")
	dbConn := db.GetDBConn()
	dailySummaryRepo := postgres.NewPostgresPowerDailySummaryRepo(dbConn)
	dailyStatusRepo := postgres.NewPostgresDeviceStatusRepo(dbConn)
	dailySummaryService := dailyas.NewDailySummaryLogic(dailySummaryRepo, dailyStatusRepo)

	slog.Info("daily summary start")
	dailySummaryService.Execute(utils.GetTimeNow())
	slog.Info("daily summary end")
}
