# Local Testing CA

## Usage
CA_PASSWORD=my-ca-password ./build.sh

## Troubleshooting
### openssl configuration for Apple based systems
etc/ssl/openssl.cnf

```
[ v3_ca ]
basicConstraints = critical,CA:TRUE
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
```

check key: openssl x509 -in mosquitto-broker/broker.crt -text -noout