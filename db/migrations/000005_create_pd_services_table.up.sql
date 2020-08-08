CREATE TABLE pagerduty_services
(
    id                   varchar             NOT NULL
        CONSTRAINT pagerduty_services_pk
            PRIMARY KEY,
    summary              varchar             NOT NULL,
    name                 varchar             NOT NULL,
    description          varchar             NOT NULL,
    created_at           timestamp           NOT NULL,
    last_incident        timestamp,
    escalation_policy_id varchar             NOT NULL,
    team_ids             character varying[] NOT NULL
);

ALTER TABLE pagerduty_services
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX pagerduty_services_id_uindex
    ON pagerduty_services (id);
