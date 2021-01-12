CREATE DATABASE evcharger
    WITH 
    OWNER = mypguser
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

CREATE TABLE public.ev_readings
(
    date date,
    charger_name character(8),
    current_output_power double precision,
    current_output_current double precision,
    current_total_energy double precision,
    connected_vehicle integer,
    current_energy_price double precision,
    current_charging_duration integer,
    current_charging_energy double precision,
    charging_state character(1)
);

ALTER TABLE public.ev_readings
    OWNER to mypguser;

SELECT create_hypertable('ev_readings', 'date')