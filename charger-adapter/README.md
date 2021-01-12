## Build

For raspberry pi:
`env GOOS=linux GOARCH=arm GOARM=5 go build`

## Modes

### Emit
Lokales Testsetup:
```
 report --charger-ip=127.0.0.1 
  --charger-pin-1=0000 
  --broker-uri=ssl://127.0.0.1:8883 
  --charger-name=CHARGER1
  --client-certificate=../charger-ca/charger1/client.crt 
  --client-key=../charger-ca/charger1/client.key 
  --trusted-ca=../charger-ca/rootca/ca-crt.key 
  --verbose
```

## MQTT

Best Practices:
- One Topic per Value (data Retain)

### Topics


examples:

* chargers/SG1000/
    * TotalEnergyConsumption 
        
Payload:
```
{
    val : 123.45 // value
    unit: kWh // unit of sampling
    ts : 2020-01-01T11:11:11Z // last sampling time
}
```

    
    

## Database

```postgresql
-- Table: public.ev_readings

-- DROP TABLE public.ev_readings;

CREATE TABLE public.ev_readings
(
    date date NOT NULL,
    charger_name character(8) COLLATE pg_catalog."default",
    current_output_power double precision,
    current_output_current double precision,
    current_total_energy double precision,
    connected_vehicle integer,
    current_energy_price double precision,
    current_charging_duration integer,
    current_charging_energy double precision,
    charging_state character(1) COLLATE pg_catalog."default"
)

TABLESPACE pg_default;

ALTER TABLE public.ev_readings
    OWNER to mypguser;
-- Index: ev_readings_date_idx

-- DROP INDEX public.ev_readings_date_idx;

CREATE INDEX ev_readings_date_idx
    ON public.ev_readings USING btree
    (date DESC NULLS FIRST)
    TABLESPACE pg_default;

-- Trigger: ts_insert_blocker

-- DROP TRIGGER ts_insert_blocker ON public.ev_readings;

CREATE TRIGGER ts_insert_blocker
    BEFORE INSERT
    ON public.ev_readings
    FOR EACH ROW
    EXECUTE PROCEDURE _timescaledb_internal.insert_blocker();
```


or:
```postgresql
CREATE TYPE OperationState as enum(
    'unknown',
    'off',
    'idle',
    'charging',
    'scheduled_downtime',
    'unscheduled_downtime'
);

CREATE TABLE public.ev_readings
(
    time TIMESTAMPTZ NOT NULL,
    OperationState OperationState,
    TotalEnergyConsumption double precision,
    PowerOutput double precision,
    OutputCurrent double precision,
    MaximumOutputCurrent double precision,            
	ConnectedVehicle varchar(32)
);

SELECT create_hypertable('public.ev_readings', 'time')
```