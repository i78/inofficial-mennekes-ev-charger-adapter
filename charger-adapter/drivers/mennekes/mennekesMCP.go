package mennekes

import (
	"encoding/json"
	"github.com/codecyclist/ev-charger-adapter/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MennekesMCP struct {
	cfg             MennekesConfiguration
	changerEndpoint string
}

type MennekesConfiguration struct {
	ChargerHost string
	ChargerPin1 string
}

type MennekesChargeDataPDO struct {
	ChgState    string `json:"ChgState"`
	Price       int    `json:"Price"`
	ChgDuration int    `json:"ChgDuration"`
	ChgNrg      int    `json:"ChgNrg"`
	ActPwr      int    `json:"ActPwr"`
	ActCurr     int    `json:"ActCurr"`
}

func NewMennekesMCP(cfg MennekesConfiguration) *MennekesMCP {
	return &MennekesMCP{
		cfg:             cfg,
		changerEndpoint: cfg.ChargerHost + "/MHCP/1.0/ChargeData?DevKey=" + cfg.ChargerPin1}
}

func (m *MennekesMCP) GetCurrentChargerStatus() (status models.ChargerStatus, err error) {
	resp, errHttp := http.Get(m.changerEndpoint)

	if errHttp != nil {
		log.WithFields(log.Fields{
			"error": errHttp,
		}).Warn("Unable to receive new Reading from charger")

		err = errHttp
		return
	}

	var chargeDataPDO MennekesChargeDataPDO

	if jsonDecodeError := json.NewDecoder(resp.Body).Decode(&chargeDataPDO); jsonDecodeError == nil {
		log.WithFields(log.Fields{
			"reading": chargeDataPDO,
		}).Trace("Received new Reading from Charger")

		status = models.ChargerStatus{
			OperationState:         mapMennekesState(chargeDataPDO.ChgState),
			TotalEnergyConsumption: float64(chargeDataPDO.ChgNrg),
			PowerOutput:            float64(chargeDataPDO.ActPwr),
			OutputCurrent:          float64(chargeDataPDO.ActCurr),
		}
	} else {
		log.WithFields(log.Fields{
			"jsonDecodeError": jsonDecodeError,
		}).Error("Unable to read data from Mennekes Charger")
		err = jsonDecodeError
	}

	return
}

func mapMennekesState(state string) models.OperationState {
	switch state {
	case "Idle":
		return models.IDLE
	case "Charging":
		return models.CHARGING
	default:
		return models.UNKNOWN
	}

}
