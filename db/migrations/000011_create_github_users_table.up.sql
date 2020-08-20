CREATE TABLE github_users
(
    id         integer NOT NULL
        CONSTRAINT github_users_pk
            PRIMARY KEY,
    login      varchar NOT NULL,
    avatar_url varchar NOT NULL,
    type       varchar NOT NULL,
    site_admin boolean NOT NULL
);

ALTER TABLE github_users
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX github_users_id_uindex
    ON github_users (id);

CREATE INDEX github_users_login_index
    ON github_users (login);
