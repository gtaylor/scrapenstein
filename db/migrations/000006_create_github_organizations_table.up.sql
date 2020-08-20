CREATE TABLE github_organizations
(
    id         integer NOT NULL
        CONSTRAINT github_organizations_pk
            PRIMARY KEY,
    login      varchar NOT NULL,
    avatar_url varchar NOT NULL
);

ALTER TABLE github_organizations
    OWNER TO scrapenstein;

CREATE UNIQUE INDEX github_organizations_id_uindex
    ON github_organizations (id);

CREATE UNIQUE INDEX github_organizations_login_uindex
    ON github_organizations (login);
