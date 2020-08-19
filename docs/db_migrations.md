For a more detailed walkthrough on DB migrations, see the [upstream docs](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md) for golang-migrate/migrate. For a tldr, see below.

# Installing golang-migrate

```
go get -u github.com/golang-migrate/migrate
```

# Creating a new migration

```shell script
migrate create -ext sql -dir db/migrations -seq your_migration_name
```

Replace `your_migration_name` with a descriptive name for the new migration.

# Migrating up or down

```shell script
# Migrate up
migrate -database postgresql://scrapenstein:scrapenstein@localhost:5432?sslmode=disable \
  -path db/migrations up
# Migrate down
migrate -database postgresql://scrapenstein:scrapenstein@localhost:5432?sslmode=disable \
  -path db/migrations down
```
