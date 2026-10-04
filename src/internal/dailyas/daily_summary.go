package dailyas

import (
	"fmt"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/dto"
	"github.com/yamaki-87/switchbot-app/src/internal/repo"
	"github.com/yamaki-87/switchbot-app/src/internal/status"
	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

type DailySummaryLogic struct {
	powerDailySummaryRepo repo.PowerDailySummaryAndNotifyRepo
	deviceStatusRepo      repo.DeviceStatusRepo
}

func NewDailySummaryLogic(powerDailySummaryRepo repo.PowerDailySummaryAndNotifyRepo, deviceStatusRepo repo.DeviceStatusRepo) *DailySummaryLogic {
	return &DailySummaryLogic{
		powerDailySummaryRepo: powerDailySummaryRepo,
		deviceStatusRepo:      deviceStatusRepo,
	}
}

func (l *DailySummaryLogic) Execute(endDate time.Time) error {
	maxDate, err := l.powerDailySummaryRepo.FindMaxProcessDate()
	if err != nil {
		return err
	}
	if maxDate == nil {
		return fmt.Errorf("maxDate is nil")
	}

	if err := l.registerDeviceSummaryAndNotify(maxDate, endDate); err != nil {
		return err
	}

	return nil
}

func (l *DailySummaryLogic) registerDeviceSummaryAndNotify(maxDate *time.Time, endDate time.Time) error {
	startDate := utils.AddDays(*maxDate, 1)

	processDates := utils.ResolveProcessDate(startDate, endDate)
	for _, processDate := range processDates {
		endDate = utils.AddDays(processDate, 1)

		deviceSummary, err := l.deviceStatusRepo.SelectDailySummary(dto.SelectDailySummaryIn{
			StartDate: processDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}

		deviceSummaries, notifies := l.deviceSummaryToInsertAndNotify(deviceSummary, processDate)

		if err := l.powerDailySummaryRepo.InsertBatch(deviceSummaries, notifies); err != nil {
			return err
		}
	}

	return nil
}

func (l *DailySummaryLogic) deviceSummaryToInsertAndNotify(deviceSummary []repo.SelectDailySummary, processDate time.Time) ([]repo.PowerDailySummary, []repo.Notify) {
	var deviceSummaries = make([]repo.PowerDailySummary, 0, len(deviceSummary))
	var notifies = make([]repo.Notify, 0, len(deviceSummary))
	for _, summary := range deviceSummary {
		deviceSummaries = append(deviceSummaries, repo.PowerDailySummary{
			DeviceID:    summary.DeviceID,
			ProcessDate: processDate,
			SumKwh:      summary.SumKwh,
		})

		notifies = append(notifies, repo.Notify{
			DeviceId:  summary.DeviceID,
			Status:    status.NotifyPending,
			NotifyMsg: fmt.Sprintf("Daily summary %s: %.2f kWh", summary.DeviceID, summary.SumKwh),
		})
	}

	return deviceSummaries, notifies
}
