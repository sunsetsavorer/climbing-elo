CREATE TABLE climbers (
    id           TEXT        PRIMARY KEY,
    telegram_id  BIGINT      NOT NULL UNIQUE,
    nickname     TEXT        NOT NULL,
    rating       INTEGER     NOT NULL DEFAULT 1000,
    ascents_count INTEGER    NOT NULL DEFAULT 0,
    hardest_grade   TEXT,
    created_at   TIMESTAMPTZ NOT NULL
);
