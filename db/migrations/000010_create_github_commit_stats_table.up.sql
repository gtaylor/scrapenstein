CREATE TABLE github_commit_stats
(
    repo_id   integer NOT NULL,
    sha       varchar NOT NULL,
    additions integer NOT NULL,
    deletions integer NOT NULL,
    total     integer NOT NULL,
    CONSTRAINT github_commits_changes_pk
        UNIQUE (repo_id, sha)
);

ALTER TABLE github_commit_stats
    OWNER TO scrapenstein;

CREATE INDEX github_commits_changes_repo_id_sha_index
    ON github_commit_stats (repo_id, sha);
