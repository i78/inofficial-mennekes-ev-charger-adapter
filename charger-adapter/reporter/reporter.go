package reporter

import (
	"context"
	"github.com/codecyclist/ev-charger-adapter/models"
	log "github.com/sirupsen/logrus"
	"time"
)

const ChargerStatusRefreshIntervalSeconds = 5
const RequestTimeout = 4000

type Reporter struct {
	ReporterConfiguration
	ChargerDriver
	context.Context
	Emitter
}

func NewReporter(cfg ReporterConfiguration, drv ChargerDriver, ctx context.Context, emitter Emitter) (r *Reporter) {
	return &Reporter{
		cfg,
		drv,
		ctx,
		emitter,
	}
}

type ChargerDriver interface {
	GetCurrentChargerStatus() (status models.ChargerStatus, err error)
}

type Emitter interface {
	EmitChargerStatus(sampletime time.Time, status models.ChargerStatus)
}

func (r *Reporter) Run() {
	log.WithFields(log.Fields{
		"ChargerName": r.ChargerName,
		"ChargerType": r.ChargerType,
	}).Info("Starting Reporter")

	for {
		select {
		case <-time.After(ChargerStatusRefreshIntervalSeconds * time.Second):
			log.Trace("Refreshing status")
			currentStatus, _ := r.ChargerDriver.GetCurrentChargerStatus()
			log.WithFields(log.Fields{
				"newChargerStatus": currentStatus,
			}).Info()
			r.Emitter.EmitChargerStatus(time.Now(), currentStatus)

		case <-r.Context.Done():
			log.Info("Terminating Reporter")
			return
		}

	}

}
