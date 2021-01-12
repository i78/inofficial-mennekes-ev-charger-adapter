# Electric Vehicle Charger Adapter for Menneken Chargers

## Abstract

This project demonstrates how to capture charger state data from Menneken EV Chargers using the same network interface
as the mobile application. The given solution pushes the data to an MQTT topic, so another instance can persist it in a
database, such as the time-series database TimescaleDB.

## Mandatory Legal Disclaimer
This is not an official Menneken project, nor I am associated with Menneken in any respect. I am just a customer who wants to be in charge of the data generated by the EV charger I paid for.

Use this at your own risk. I am not responsible for any damage or undesired operation caused by the use of this program. Do not use this application if you don't agree with that.

If you are associated with Menneken and are not comfortable with this application, please contact me so I can alter this project according to your advise. Thank you!

## This Repository

| Folder | Description |
| --- | --- |
| charger-adapter| The actual adapter and persistor |
| backend | A local docker-compose stack for spinning up a TimescaleDB and Grafana instance for analysis |
| charger-ca | A convenience script to start an own certificate authority as a quick seed for Mosquitto MQTT certificate-based authentication |
| mock-charger| A tiny node based mock HTTP server, mimicing a charger for local demonstration purposes.|| charger-adapter| The main project


## Usage

### Commands

Listen to charger telegrams:

```bash
(app) mosquitto_sub -v -h localhost -t chargers/statusUpdate/SG1
```

Persist incoming telegrams to database:

```bash
(app) mosquitto_sub -v -h localhost -t chargers/statusUpdate/SG1
```

## Code of Conduct
This project adheres to the [Code of Merit](https://codeofmerit.org/code/).

## References
Thanks to:
- http://www.steves-internet-guide.com/mosquitto-tls/
- https://github.com/mchestr/Secure-MQTT-Docker