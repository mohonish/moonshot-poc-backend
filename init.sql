CREATE database iss_data;
\c iss_data;
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE TABLE location ( time TIMESTAMPTZ NOT NULL, latitude DOUBLE PRECISION NULL, longitude DOUBLE PRECISION NULL );
SELECT create_hypertable('location', 'time');