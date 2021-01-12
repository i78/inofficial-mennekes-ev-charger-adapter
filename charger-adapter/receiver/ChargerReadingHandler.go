package receiver

import "github.com/codecyclist/ev-charger-adapter/models"

type ChargerReadingHandler interface {
	Handle(readingName string, reading models.ChargerValueEnvelope) error
}
