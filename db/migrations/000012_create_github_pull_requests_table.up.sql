CREATE TABLE github_pull_requests
(
    id                     integer             NOT NULL
        CONSTRAINT github_pull_requests_pk
            PRIMARY KEY,
    repo_id                integer             NOT NULL,
    number                 integer             NOT NULL,
    state                  varchar             NOT NULL,
    locked                 boolean             NOT NULL,
    title                  varchar             NOT NULL,
    user_id                integer             NOT NULL,
    created_at             timestamp           NOT NULL,
    updated_at             timestamp           NOT NULL,
    closed_at              timestamp,
    merged_at              timestamp,
    labels                 character varying[] NOT NULL,
    draft                  boolean             NOT NULL,
    merged                 boolean             NOT NULL,
    mergeable              boolean             NOT NULL,
    merged_by_id           integer,
    rebaseable             boolean             NOT NULL,
    comments               integer             NOT NULL,
    review_comments        integer             NOT NULL,
    commits                integer             NOT NULL,
    additions              integer             NOT NULL,
    deletions              integer             NOT NULL,
    changed_files          integer             NOT NULL,
    assignee_ids           integer[]           NOT NULL,
    requested_reviewer_ids integer[]           NOT NULL
);

ALTER TABLE github_pull_requests
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX github_pull_requests_id_uindex
    ON github_pull_requests (id);
