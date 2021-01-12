package mqtt

import (
	"github.com/codecyclist/ev-charger-adapter/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmitter(t *testing.T) {
	emitter := MQTTEmitter{
		client:                  nil,
		config:                  Config{},
		previouslyEmittedStatus: models.ChargerStatus{},
	}

	currentStatus := models.ChargerStatus{
		OperationStatus:        models.IDLE,
		TotalEnergyConsumption: 1.0,
		PowerOutput:            1.0,
		OutputCurrent:          1.0,
		MaximumOutputCurrent:   1.0,
	}

	t.Run("should report changed value when one changed", func(t *testing.T) {
		nextStatus := models.ChargerStatus{
			OperationStatus:        models.IDLE,
			TotalEnergyConsumption: 2.0,
			PowerOutput:            1.0,
			OutputCurrent:          1.0,
			MaximumOutputCurrent:   1.0,
		}

		changes := emitter.deriveMessagesFromChangedValues(currentStatus, nextStatus)

		assert.Equal(t, 2.0, changes[0].Payload.Value)
	})

	t.Run("should report changed non-float when one changed", func(t *testing.T) {
		nextStatus := models.ChargerStatus{
			OperationStatus:        models.CHARGING,
			TotalEnergyConsumption: 1.0,
			PowerOutput:            1.0,
			OutputCurrent:          1.0,
			MaximumOutputCurrent:   1.0,
		}

		changes := emitter.deriveMessagesFromChangedValues(currentStatus, nextStatus)

		assert.Equal(t,1, len(changes))
		assert.Equal(t, models.CHARGING, changes[0].Payload.Value)
	})

	t.Run("should report changed values when multiple changed", func(t *testing.T) {
		nextStatus := models.ChargerStatus{
			OperationStatus:        models.IDLE,
			TotalEnergyConsumption: 2.0,
			PowerOutput:            2.0,
			OutputCurrent:          2.0,
			MaximumOutputCurrent:   2.0,
		}

		changes := emitter.deriveMessagesFromChangedValues(currentStatus, nextStatus)

		assert.Equal(t,4, len(changes))
		for _, change := range changes {
			assert.Equal(t, 2.0, change.Payload.Value)
		}
	})

	t.Run("should report operation status as string when changed", func(t *testing.T) {
		nextStatus := models.ChargerStatus{
			OperationStatus:        models.CHARGING,
			TotalEnergyConsumption: 1.0,
			PowerOutput:            1.0,
			OutputCurrent:          1.0,
			MaximumOutputCurrent:   1.0,
		}

		changes := emitter.deriveMessagesFromChangedValues(currentStatus, nextStatus)

		assert.Equal(t, models.CHARGING, changes[0].Payload.Value)

	})

}
