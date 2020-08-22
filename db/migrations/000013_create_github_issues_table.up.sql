CREATE TABLE github_issues
(
    id           integer             NOT NULL
        CONSTRAINT github_issues_pk
            PRIMARY KEY,
    repo_id      integer             NOT NULL,
    number       integer             NOT NULL,
    state        varchar             NOT NULL,
    locked       boolean             NOT NULL,
    title        varchar             NOT NULL,
    user_id      integer             NOT NULL,
    created_at   timestamp           NOT NULL,
    updated_at   timestamp           NOT NULL,
    closed_at    timestamp,
    labels       character varying[] NOT NULL,
    comments     integer             NOT NULL,
    assignee_ids integer[]           NOT NULL
);

ALTER TABLE github_issues
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX github_issues_id_uindex
    ON github_issues (id);
