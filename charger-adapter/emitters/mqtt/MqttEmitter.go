// This package contains an implementation of the Emitter interface which uses MQTT to report the charger status
package mqtt


import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/codecyclist/ev-charger-adapter/models"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"reflect"
	"time"
)

// MQTTEmitter implements the Emitter Interface towards an MQTT client
type MQTTEmitter struct {
	client MQTT.Client
	config Config
	// previouslyEmittedStatus contains a copy of the previously emitted status for change detection
	previouslyEmittedStatus models.ChargerStatus
}

func NewMQTTEmitter(cfg Config) (e *MQTTEmitter) {
	// inspired by
	// https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/ssl/main.go
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(cfg.TrustedCACertificate)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(cfg.ClientCertificate, cfg.ClientKey)
	if err != nil {
		panic(err)
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.Leaf.DNSNames)

	tlsConfig := tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		// You might want to set this to false for production. :)
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}

	opts := MQTT.NewClientOptions().
		AddBroker(cfg.Broker).
		SetTLSConfig(&tlsConfig).
		SetClientID(cfg.ClientId)

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTTEmitter{client: client, config: cfg}
}

func hasValueChanged(currentField reflect.Value, nextField reflect.Value) (hasChanged bool, currentValue interface{}, newValue interface{}) {
	currentValue = currentField.Interface()
	newValue = nextField.Interface()
	hasChanged = currentValue != newValue

	return
}

func (e *MQTTEmitter) deriveMessagesFromChangedValues(currentStatus models.ChargerStatus, nextStatus models.ChargerStatus) (changeMessages []models.ReadyToSendValue) {
	current := reflect.ValueOf(&currentStatus).Elem()
	next := reflect.ValueOf(&nextStatus).Elem()
	chargerStatusType := current.Type()

	for i := 0; i < chargerStatusType.NumField(); i++ {
		fieldName := chargerStatusType.Field(i).Name
		currentField := current.FieldByName(fieldName)
		nextField := next.FieldByName(fieldName)

		if changed, currentValue, newValue := hasValueChanged(currentField, nextField); changed {
			log.WithFields(log.Fields{
				"fieldName":    fieldName,
				"currentValue": currentValue,
				"nextValue":    newValue,
			}).
				Info("Change detected")

			changeMessages = append(changeMessages, models.ReadyToSendValue{
				Payload: models.ChargerValueEnvelope{
					SamplingTime: time.Now(),
					Value:        newValue,
					Unit:         "any",
				},
				Topic: e.config.Topic + fieldName,
			})
		}
	}

	return
}

func (e *MQTTEmitter) EmitChargerStatus(sampletime time.Time, nextStatus models.ChargerStatus) {
	dueMessages := e.deriveMessagesFromChangedValues(e.previouslyEmittedStatus, nextStatus)
	e.previouslyEmittedStatus = nextStatus

	for _, message := range dueMessages {
		payloadJson, _ := json.Marshal(message.Payload)

		log.WithFields(log.Fields{
			"sampletime": sampletime,
			"topic":      message.Topic,
			"value":      message.Payload.Value,
		}).Debug("Emitting new Status")

		e.client.Publish(message.Topic, 0, true, payloadJson)
	}
}
