package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/yamaki-87/switchbot-app/src/internal/app"
	"github.com/yamaki-87/switchbot-app/src/internal/collector"
	"github.com/yamaki-87/switchbot-app/src/internal/config"
	"github.com/yamaki-87/switchbot-app/src/internal/db"
	"github.com/yamaki-87/switchbot-app/src/internal/devicemaster"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func main() {
	_ = godotenv.Load()

	err := config.LoadSettings(os.Getenv("SETTING_PATH"))
	utils.FatalIfErr(err, "load settings failed")

	settings := config.GetSettings()
	client := &http.Client{
		Timeout: time.Duration(settings.Config.App.IntervalTimeout) * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = db.DbInit(ctx, settings.Config.Database.URL)
	utils.FatalIfErr(err, "db init failed")

	dbConn := db.GetDBConn()
	deviceMasterRepo := repo.NewPostgresDeviceMasterRepo(dbConn)
	deviceMasterManager, err := devicemaster.NewDeviceMasterManager(deviceMasterRepo)
	utils.FatalIfErr(err, "new device master manager failed")

	deviceStatusRepo := repo.NewPostgresDeviceStatusRepo(dbConn)
	deviceMaster := deviceMasterManager.GetDeviceMaster()

	interval := deviceMaster[0].PollingIntervalSec
	colletorLogic := collector.NewCollectorLogic(deviceStatusRepo)
	collectorApp := app.NewCollectorApp(
		client,
		colletorLogic,
		deviceMaster,
		settings.Secret.Token,
		settings.Secret.Secret,
	)

	collectorApp.Run(ctx, time.Duration(interval)*time.Second)
}
