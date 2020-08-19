CREATE TABLE github_commits
(
    repo_id                integer             NOT NULL,
    sha                    varchar             NOT NULL,
    author_id              integer             NOT NULL,
    committer_id           integer             NOT NULL,
    parents_sha            character varying[] NOT NULL,
    commit_author_name     varchar             NOT NULL,
    commit_author_email    varchar             NOT NULL,
    commit_author_date     timestamp           NOT NULL,
    commit_committer_name  varchar             NOT NULL,
    commit_committer_email varchar             NOT NULL,
    commit_committer_date  timestamp           NOT NULL,
    message                text                NOT NULL,
    tree_sha               varchar             NOT NULL,
    verification_verified  boolean             NOT NULL,
    verification_reason    varchar             NOT NULL,
    CONSTRAINT github_commits_pk
        PRIMARY KEY (repo_id, sha)
);

ALTER TABLE github_commits
    OWNER TO scrapenstein;

CREATE INDEX github_commits_repo_id_sha_index
    ON github_commits (repo_id, sha);
