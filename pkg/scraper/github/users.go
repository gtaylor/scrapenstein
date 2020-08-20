package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeUsersOptions struct{}

func ScrapeUsers(dbConn *pgx.Conn, ghOptions GitHubOptions, options ScrapeUsersOptions) (int, error) {
	client, err := newGHClient(ghOptions)
	if err != nil {
		return 0, err
	}
	listAllOpts := github.UserListOptions{}
	totalUsers := 0
	for {
		users, _, err := client.Users.ListAll(context.Background(), &listAllOpts)
		if err != nil {
			return totalUsers, err
		}
		if len(users) == 0 {
			break
		}
		for _, user := range users {
			if err := storeUser(dbConn, user); err != nil {
				return totalUsers, err
			}
			totalUsers += 1
		}
		lastUser := users[len(users)-1]
		listAllOpts.Since = *lastUser.ID
	}
	return totalUsers, nil
}

const storeUserQuery = `
	INSERT INTO github_users (
		id, login, avatar_url, type, site_admin
	) VALUES(
		$1, $2, $3, $4, $5
	) ON CONFLICT (id)
		DO UPDATE SET
			id=excluded.id,
			login=excluded.login,
			avatar_url=excluded.avatar_url,
			type=excluded.type,
			site_admin=excluded.site_admin`

func storeUser(dbConn *pgx.Conn, user *github.User) error {

	_, err := dbConn.Exec(
		context.Background(), storeUserQuery,
		user.GetID(), user.GetLogin(), user.GetAvatarURL(), user.GetType(),
		user.GetSiteAdmin())
	return err
}
