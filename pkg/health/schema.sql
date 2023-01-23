DROP TYPE IF EXISTS workout_statistics_t;
CREATE TYPE workout_statistics_t AS (
    type       CHARACTER VARYING,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date   TIMESTAMP WITH TIME ZONE,
    average    CHARACTER VARYING,
    minimum    CHARACTER VARYING,
    maximum    CHARACTER VARYING,
    sum        DECIMAL,
    unit       CHARACTER VARYING
);

DROP TYPE IF EXISTS metadata_t;
CREATE TYPE metadata_t AS (
    key   CHARACTER VARYING,
    value CHARACTER VARYING
);

DROP TABLE IF EXISTS records;
CREATE TABLE IF NOT EXISTS records (
    id             SERIAL PRIMARY KEY,
    type           CHARACTER VARYING NOT NULL,
    unit           CHARACTER VARYING,
    value          CHARACTER VARYING,
    source_name    CHARACTER VARYING NOT NULL,
    source_version CHARACTER VARYING,
    device         CHARACTER VARYING,
    creation_date  TIMESTAMP WITH TIME ZONE,
    start_date     TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date       TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata       JSONB,
    hrv            JSONB
);

DROP TABLE IF EXISTS workouts;
CREATE TABLE IF NOT EXISTS workouts (
    id                       SERIAL PRIMARY KEY,
    workout_activity_type    CHARACTER VARYING NOT NULL,
    duration                 DECIMAL,
    duration_unit            CHARACTER VARYING,
    total_distance           CHARACTER VARYING,
    total_distance_unit      CHARACTER VARYING,
    total_energy_burned      CHARACTER VARYING,
    total_energy_burned_unit CHARACTER VARYING,
    source_name              CHARACTER VARYING NOT NULL,
    source_version           CHARACTER VARYING,
    device                   CHARACTER VARYING,
    creation_date            TIMESTAMP WITH TIME ZONE,
    start_date               TIMESTAMP WITH TIME ZONE,
    end_date                 TIMESTAMP WITH TIME ZONE,

    metadata           JSONB,  -- array of metadata_t
    workout_events     JSONB,
    workout_routes     JSONB,
    workout_statistics JSONB   -- array of workout_statistics_t
);

DROP TABLE IF EXISTS activity_summaries;
CREATE TABLE IF NOT EXISTS activity_summaries (
    date_components           CHARACTER VARYING,
    active_energy_burned      CHARACTER VARYING,
    active_energy_burned_goal CHARACTER VARYING,
    active_energy_burned_unit CHARACTER VARYING,
    apple_move_time           CHARACTER VARYING,
    apple_move_time_goal      CHARACTER VARYING,
    apple_exercise_time       CHARACTER VARYING,
    apple_exercise_time_goal  CHARACTER VARYING,
    apple_stand_hours         CHARACTER VARYING,
    apple_stand_hours_goal    CHARACTER VARYING
);

DROP TABLE IF EXISTS clinical_records;
CREATE TABLE IF NOT EXISTS clinical_records (
    type               CHARACTER VARYING,
    identifier         CHARACTER VARYING,
    source_name        CHARACTER VARYING,
    source_url         CHARACTER VARYING,
    fhir_version       CHARACTER VARYING,
    received_date      CHARACTER VARYING,
    resource_file_path CHARACTER VARYING
);

DROP TABLE IF EXISTS audiograms;
CREATE TABLE IF NOT EXISTS audiograms (
    type          CHARACTER VARYING NOT NULL,
    sourceName    CHARACTER VARYING NOT NULL,
    sourceVersion CHARACTER VARYING,
    device        CHARACTER VARYING,
    creationDate  TIMESTAMP WITH TIME ZONE,
    startDate     TIMESTAMP WITH TIME ZONE NOT NULL,
    endDate       TIMESTAMP WITH TIME ZONE NOT NULL,

    metadata           JSONB,
    sensitivity_points JSONB
);

DROP TABLE IF EXISTS vision_prescriptions;
CREATE TABLE IF NOT EXISTS vision_prescriptions (
    type             CHARACTER VARYING NOT NULL,
    dateIssued       CHARACTER VARYING NOT NULL,
    expirationDate   CHARACTER VARYING,
    brand            CHARACTER VARYING,

    metadata    JSONB,
    right_eye   JSONB,
    left_eye    JSONB,
    attachments JSONB
);