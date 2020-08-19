CREATE TABLE github_commit_files
(
    repo_id   integer NOT NULL,
    sha       varchar NOT NULL,
    filename  varchar NOT NULL,
    additions integer NOT NULL,
    deletions integer NOT NULL,
    changes   integer NOT NULL,
    status    varchar NOT NULL,
    CONSTRAINT github_commits_files_pk
        UNIQUE (repo_id, sha, filename)
);

ALTER TABLE github_commit_files
    OWNER TO scrapenstein;

CREATE INDEX github_commits_files_repo_id_sha_index
    ON github_commit_files (repo_id, sha);
