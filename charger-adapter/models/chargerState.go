package models

import "time"

type OperationState string

const (
	UNKNOWN              OperationState = "unknown"
	OFF                  OperationState = "off"
	IDLE                 OperationState = "idle"
	CHARGING             OperationState = "charging"
	SCHEDULED_DOWNTIME   OperationState = "scheduled_downtime"
	UNSCHEDULED_DOWNTIME OperationState = "unscheduled_downtime"
)

type ChargerStatus struct {
	OperationState        OperationState
	EnergyPrice            float64
	TotalEnergyConsumption float64
	PowerOutput            float64
	OutputCurrent          float64
	MaximumOutputCurrent   float64
}

type ChargerValueEnvelope struct {
	SamplingTime time.Time
	Value        interface{}
	Unit         string
}

type ReadyToSendValue struct {
	Payload ChargerValueEnvelope
	Device  string
	Topic   string
}
