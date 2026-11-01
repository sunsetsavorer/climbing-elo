CREATE TABLE ascents (
    id           TEXT        PRIMARY KEY,
    climber_id   TEXT        NOT NULL REFERENCES climbers(id) ON DELETE CASCADE,
    grade        TEXT        NOT NULL,
    outcome      TEXT        NOT NULL,
    delta        INTEGER     NOT NULL,
    rating_after INTEGER     NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL
);

CREATE INDEX ascents_climber_id_idx ON ascents (climber_id);
