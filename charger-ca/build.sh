# MQTT Broker test CA and certs for development purposes
RSA_KEY_SIZE=4096

## CA
CA_PRIV_KEY_FILE=rootca/ca-key.pem
CA_ROOT_CERTIFICATE_FILE=rootca/ca-crt.key

### Create Certificate Authority
openssl genpkey -aes-256-cbc \
    -pass pass:$CA_PASSWORD \
    -algorithm RSA \
    -out $CA_PRIV_KEY_FILE \
    -pkeyopt rsa_keygen_bits:$RSA_KEY_SIZE

### Create Root Certificate
openssl req -x509 -new -nodes -extensions v3_ca \
    -key $CA_PRIV_KEY_FILE \
    -days 1024 \
    -out $CA_ROOT_CERTIFICATE_FILE \
    -sha512 \
    -passin pass:$CA_PASSWORD \
    -subj "/C=DE/ST=NRW/L=Solingen/O=chargers/CN=evcharger.internal"


## Mosquitto Broker
MOSQUITTO_BROKER_PRIV_KEY_FILE=mosquitto-broker/broker.key
MOSQUITTO_BROKER_PRIV_KEY_FILE_CSR=mosquitto-broker/broker.pem
MOSQUITTO_BROKER_PRIV_KEY_PASSPHRASE='passphrase'
MOSQUITTO_BROKER_CERTIFICATE=mosquitto-broker/broker.crt

### Generate Key and CSR
openssl req -new -nodes -subj "/C=DE/ST=NRW/L=Solingen/O=development_certificate/CN=127.0.0.1" \
    -newkey rsa:$RSA_KEY_SIZE \
    -reqexts SAN -config san.cfg \
    -keyout $MOSQUITTO_BROKER_PRIV_KEY_FILE \
    -out $MOSQUITTO_BROKER_PRIV_KEY_FILE_CSR


### Sign Broker Certificate
openssl x509 -req \
    -in $MOSQUITTO_BROKER_PRIV_KEY_FILE_CSR  \
    -CA $CA_ROOT_CERTIFICATE_FILE \
    -CAkey $CA_PRIV_KEY_FILE \
    -CAcreateserial \
    -out $MOSQUITTO_BROKER_CERTIFICATE \
    -passin pass:$CA_PASSWORD  \
    -days 365 \
    -sha512

## Charger1 Client Certificate
CHARGER1_CLIENT_KEY=charger1/client.key
CHARGER1_CLIENT_CSR=charger1/client.csr
CHARGER1_CLIENT_CERT=charger1/client.crt

### Generate Key and CSR

openssl req -new -nodes -subj "/C=DE/ST=NRW/L=Solingen/O=development_certificate/CN=127.0.0.1" \
    -newkey rsa:$RSA_KEY_SIZE \
    -reqexts SAN -config san.cfg \
    -keyout $CHARGER1_CLIENT_KEY \
    -out $CHARGER1_CLIENT_CSR

### Sign Key
openssl x509 -req \
    -in $CHARGER1_CLIENT_CSR  \
    -CA $CA_ROOT_CERTIFICATE_FILE \
    -CAkey $CA_PRIV_KEY_FILE \
    -CAcreateserial \
    -out $CHARGER1_CLIENT_CERT \
    -passin pass:$CA_PASSWORD  \
    -days 365 -sha512

## Install certificates
cp $CA_ROOT_CERTIFICATE_FILE ../backend/mosquitto/config/certs/rootCA.crt
cp $MOSQUITTO_BROKER_PRIV_KEY_FILE ../backend/mosquitto/config/certs/server.key
cp $MOSQUITTO_BROKER_CERTIFICATE ../backend/mosquitto/config/certs/server.crt