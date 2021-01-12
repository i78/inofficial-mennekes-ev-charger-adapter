/*
For usage, please see README.md
 */
package main

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/codecyclist/ev-charger-adapter/drivers/mennekes"
	mqtt2 "github.com/codecyclist/ev-charger-adapter/emitters/mqtt"
	persistor "github.com/codecyclist/ev-charger-adapter/persistor"
	"github.com/codecyclist/ev-charger-adapter/receiver"
	"github.com/codecyclist/ev-charger-adapter/reporter"
	log "github.com/sirupsen/logrus"
)

type verboseFlag bool

var cli struct {
	Report  reportCommand  `cmd help:"Run status report"`
	Receive receiveCommand `cmd help:"Run status report"`
	Verbose verboseFlag    `help:"Enable debug logging."`
}

type reportCommand struct {
	ChargerName       string
	BrokerUri         string
	ClientCertificate string
	ClientKey         string
	TrustedCA         string
	ChargerType       string
	ChargerIP         string
	ChargerPin1       string
}

type receiveCommand struct {
	ChargerName       string
	BrokerUri         string
	ClientCertificate string
	ClientKey         string
	TrustedCA         string
}

func (v verboseFlag) BeforeApply() error {
	log.SetLevel(log.TraceLevel)
	return nil
}

func (sv *reportCommand) Run(ctx *kong.Context) error {
	log.WithFields(log.Fields{
		"chargerIp": sv.ChargerIP,
	}).Info("Starting")

	cancel := context.Background()

	chargerDriver := mennekes.NewMennekesMCP(mennekes.MennekesConfiguration{
		ChargerHost: "http://" + sv.ChargerIP + ":25000",
		ChargerPin1: sv.ChargerPin1,
	})

	emitter := mqtt2.NewMQTTEmitter(mqtt2.Config{
		Broker:               sv.BrokerUri,
		ClientId:             "chargers_" + sv.ChargerName,
		Topic:                "chargers/" + sv.ChargerName + "/",
		ClientCertificate:    sv.ClientCertificate,
		ClientKey:            sv.ClientKey,
		TrustedCACertificate: sv.TrustedCA,
	})

	reporter := reporter.NewReporter(reporter.ReporterConfiguration{
		ChargerName:       sv.ChargerName,
		ClientCertificate: sv.ClientCertificate,
		ChargerType:       sv.ChargerType,
	}, chargerDriver, cancel, emitter)
	reporter.Run()

	return nil
}

func (sv *receiveCommand) Run(ctx *kong.Context) error {
	log.WithFields(log.Fields{}).Info("Starting")

	cancel := context.Background()

	persistor := persistor.NewPersistor(cancel)

	recv := receiver.NewMQTTSubscriber(mqtt2.Config{
		Broker:               sv.BrokerUri,
		ClientId:             "reporter",
		ClientCertificate:    sv.ClientCertificate,
		ClientKey:            sv.ClientKey,
		TrustedCACertificate: sv.TrustedCA,
	}, "chargers/"+sv.ChargerName, persistor, cancel)

	recv.Run()

	return nil
}

func initLogging() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.InfoLevel)
}

func main() {
	initLogging()
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
