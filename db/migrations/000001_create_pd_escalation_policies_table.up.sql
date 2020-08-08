CREATE TABLE pagerduty_escalation_policies
(
    id          varchar NOT NULL
        CONSTRAINT pagerduty_escalation_policies_pk PRIMARY KEY,
    name        varchar NOT NULL,
    description text    NOT NULL
);

ALTER TABLE pagerduty_escalation_policies
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX pagerduty_escalation_policies_id_uindex
    ON pagerduty_escalation_policies (id);
