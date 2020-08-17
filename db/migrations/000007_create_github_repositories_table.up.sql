CREATE TABLE github_repositories
(
    id                     integer             NOT NULL
        CONSTRAINT github_repositories_pk
            PRIMARY KEY,
    name                   varchar             NOT NULL,
    full_name              varchar             NOT NULL,
    owner_id               integer             NOT NULL,
    owner_type             varchar             NOT NULL,
    private                boolean             NOT NULL,
    description            varchar             NOT NULL,
    fork                   boolean             NOT NULL,
    url                    varchar             NOT NULL,
    forks_count            integer             NOT NULL,
    stargazers_count       integer             NOT NULL,
    watchers_count         integer             NOT NULL,
    size                   integer             NOT NULL,
    default_branch         varchar             NOT NULL,
    open_issues_count      integer             NOT NULL,
    is_template            boolean             NOT NULL,
    topics                 character varying[] NOT NULL,
    has_issues             boolean             NOT NULL,
    has_projects           boolean             NOT NULL,
    has_wiki               boolean             NOT NULL,
    has_pages              boolean             NOT NULL,
    has_downloads          boolean             NOT NULL,
    archived               boolean             NOT NULL,
    disabled               boolean             NOT NULL,
    visibility             varchar             NOT NULL,
    pushed_at              timestamp,
    created_at             timestamp           NOT NULL,
    updated_at             timestamp,
    allow_rebase_merge     boolean             NOT NULL,
    allow_squash_merge     boolean             NOT NULL,
    delete_branch_on_merge boolean             NOT NULL,
    allow_merge_commit     boolean             NOT NULL,
    subscribers_count      integer             NOT NULL,
    network_count          integer             NOT NULL,
    organization_id        integer
);

ALTER TABLE github_repositories
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX github_repositories_full_name_uindex
    ON github_repositories (full_name);

CREATE UNIQUE INDEX github_repositories_id_uindex
    ON github_repositories (id);

CREATE UNIQUE INDEX github_repositories_name_uindex
    ON github_repositories (name);
