CREATE TABLE switchbot_device (
    device_id       text PRIMARY KEY,
    device_name     text NOT NULL,
    device_type     text NOT NULL,
    hub_device_id   text,
    is_deleted      boolean NOT NULL DEFAULT false
);

CREATE TABLE switchbot_device_settings (
    device_id           text PRIMARY KEY,
    polling_interval_sec smallint NOT NULL,

    CONSTRAINT fk_switchbot_device_settings_device
        FOREIGN KEY (device_id)
        REFERENCES switchbot_device(device_id)
);

CREATE TABLE switchbot_power (
    device_id               text NOT NULL,
    create_timestamp        timestamptz NOT NULL,

    power_w                 double precision,
    voltage_v               double precision,
    current_ma              integer,
    electricity_of_day_min  smallint,
    interval_kwh            double precision,

    collection_status       smallint NOT NULL DEFAULT 0,

    PRIMARY KEY (device_id, create_timestamp),

    CONSTRAINT fk_switchbot_power_device
        FOREIGN KEY (device_id)
        REFERENCES switchbot_device(device_id)
);

CREATE INDEX idx_switchbot_power_timestamp
ON switchbot_power (create_timestamp);

CREATE TABLE switchbot_notify (
    notify_id         serial PRIMARY KEY,
    device_id         text NOT NULL,
    status            smallint NOT NULL,
    notify_msg        text,
    create_timestamp  timestamptz NOT NULL DEFAULT now(),
    update_timestamp  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE switchbot_power_daily_summary (
    device_id         text NOT NULL,
    process_date      date NOT NULL,
    sum_kwh           double precision NOT NULL,
    create_timestamp  timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY (device_id, process_date),

    CONSTRAINT fk_switchbot_power_daily_summary_device
        FOREIGN KEY (device_id)
        REFERENCES switchbot_device(device_id)
);
