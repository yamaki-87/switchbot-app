package main

import (
	"time"

	"github.com/yamaki-87/switchbot-app/src/fw"
	"github.com/yamaki-87/switchbot-app/src/internal/app"
	"github.com/yamaki-87/switchbot-app/src/internal/collector"
	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/devicemaster"
	"github.com/yamaki-87/switchbot-app/src/internal/postgres"
	"github.com/yamaki-87/switchbot-app/src/internal/switchbotapi"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func main() {
	initOut, err := fw.Init()
	utils.FatalIfErr(err, "init failed")
	defer initOut.Stop()
	gateway := switchbotapi.NewClient(initOut.Client, initOut.Settings.Secret.Token, initOut.Settings.Secret.Secret)

	err = db.DbInit(initOut.Ctx, initOut.Settings.Config.Database.URL)
	utils.FatalIfErr(err, "db init failed")

	dbConn := db.GetDBConn()
	deviceMasterRepo := postgres.NewPostgresDeviceMasterRepo(dbConn)
	deviceMasterManager, err := devicemaster.NewDeviceMasterManager(deviceMasterRepo)
	utils.FatalIfErr(err, "new device master manager failed")

	deviceStatusRepo := postgres.NewPostgresDeviceStatusRepo(dbConn)
	deviceMaster := deviceMasterManager.GetDeviceMaster()

	interval := deviceMaster[0].PollingIntervalSec
	colletorLogic := collector.NewCollectorLogic(deviceStatusRepo, gateway)
	collectorApp := app.NewCollectorApp(
		colletorLogic,
		deviceMaster,
	)

	collectorApp.Run(initOut.Ctx, time.Duration(interval)*time.Second)
}
