package receiver

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	mqtt2 "github.com/codecyclist/ev-charger-adapter/emitters/mqtt"
	"github.com/codecyclist/ev-charger-adapter/models"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

type MQTTSubscriber struct {
	client           MQTT.Client
	config           mqtt2.Config
	context          context.Context
	handler          ChargerReadingHandler
	chargerRootTopic string
}

func NewMQTTSubscriber(cfg mqtt2.Config, chargerRootTopic string, handler ChargerReadingHandler, ctx context.Context) *MQTTSubscriber {
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

	subscriber := &MQTTSubscriber{
		client:  client,
		config:  cfg,
		handler: handler,
		context: ctx}

	subscriber.chargerRootTopic = chargerRootTopic
	subscriber.client.Subscribe(chargerRootTopic+"/#", 0, subscriber.HandleNewMessage)

	return subscriber
}

func (m *MQTTSubscriber) HandleNewMessage(client MQTT.Client, message MQTT.Message) {
	var chargerReading models.ChargerValueEnvelope
	messageReader := bytes.NewReader(message.Payload())

	if errDecodeJson := json.NewDecoder(messageReader).Decode(&chargerReading); errDecodeJson == nil {
		readingName := strings.Replace(message.Topic(), m.chargerRootTopic+"/", "", -1)

		log.WithFields(log.Fields{
			"charger": message.Topic(),
			"msg":     message.MessageID(),
			"reading": chargerReading,
		}).Trace("Received new Reading from Charger")

		m.handler.Handle(readingName, chargerReading)

	} else {
		log.Warn("Unable to decode incoming status json telegram")
	}
}

func (r *MQTTSubscriber) Run() {
	log.Info("Starting Subscriber")

	for {
		select {

		case <-r.context.Done():
			log.Info("Terminating Persistor")
			return
		}

	}

}
