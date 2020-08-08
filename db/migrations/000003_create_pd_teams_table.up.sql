CREATE TABLE pagerduty_teams
(
    id          varchar NOT NULL
        CONSTRAINT pagerduty_teams_pk
            PRIMARY KEY,
    summary     varchar NOT NULL,
    name        varchar NOT NULL,
    description varchar NOT NULL
);

ALTER TABLE pagerduty_teams
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX pagerduty_teams_id_uindex
    ON pagerduty_teams (id);
