CREATE TABLE pagerduty_incidents
(
    id                    varchar             NOT NULL
        CONSTRAINT pagerduty_incidents_pk
            PRIMARY KEY,
    summary               varchar,
    incident_number       integer             NOT NULL,
    created_at            timestamp,
    status                varchar             NOT NULL,
    title                 varchar             NOT NULL,
    incident_key          varchar             NOT NULL,
    service_id            varchar             NOT NULL,
    last_status_change_at timestamp,
    escalation_policy_id  varchar             NOT NULL,
    team_ids              character varying[] NOT NULL,
    priority_id           varchar             NOT NULL,
    urgency               varchar
);

ALTER TABLE pagerduty_incidents
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX pagerduty_incidents_id_uindex
    ON pagerduty_incidents (id);
