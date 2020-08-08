CREATE TABLE pagerduty_priorities
(
    id          varchar NOT NULL
        CONSTRAINT pagerduty_priorities_pk
            PRIMARY KEY,
    summary     varchar NOT NULL,
    name        varchar NOT NULL,
    description varchar NOT NULL
);

ALTER TABLE pagerduty_priorities
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX pagerduty_priorities_id_uindex
    ON pagerduty_priorities (id);
