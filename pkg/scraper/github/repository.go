package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeRepositoryOptions struct {
	OrgRepoAndRepoId
}

func ScrapeRepository(dbConn *pgx.Conn, ghOptions GitHubOptions, options ScrapeRepositoryOptions) error {
	client, err := newGHClient(ghOptions)
	if err != nil {
		return err
	}
	repo, _, err := client.Repositories.Get(context.Background(), options.Owner, options.Repo)
	if err != nil {
		return err
	}
	if err := storeRepository(dbConn, repo); err != nil {
		return err
	}
	return nil
}

const storeRepositoryQuery = `
	INSERT INTO github_repositories (
		id, name, full_name, owner_id, owner_type, private, description, fork, url, forks_count, stargazers_count, 
		watchers_count, size, default_branch, open_issues_count, is_template, topics, has_issues, has_projects, 
		has_wiki, has_pages, has_downloads, archived, disabled, visibility, pushed_at, created_at, updated_at,
		allow_rebase_merge, allow_squash_merge, delete_branch_on_merge, allow_merge_commit, subscribers_count,
		network_count, organization_id)
	VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 
		$22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35
	) ON CONFLICT (id)
		DO UPDATE SET 
			id=excluded.id, 
			name=excluded.name, 
			full_name=excluded.full_name,
			owner_id=excluded.owner_id,
			owner_type=excluded.owner_type,
			private=excluded.private,
			description=excluded.description,
			fork=excluded.fork,
			url=excluded.url,
			forks_count=excluded.forks_count,
			stargazers_count=excluded.stargazers_count,
			watchers_count=excluded.watchers_count,
			size=excluded.size,
			default_branch=excluded.default_branch,
			open_issues_count=excluded.open_issues_count,
			is_template=excluded.is_template,
			topics=excluded.topics,
			has_issues=excluded.has_issues,
			archived=excluded.archived,
			disabled=excluded.disabled,
			visibility=excluded.visibility,
			pushed_at=excluded.pushed_at,
			created_at=excluded.created_at,
			updated_at=excluded.updated_at,
			allow_rebase_merge=excluded.allow_rebase_merge,
			allow_squash_merge=excluded.allow_squash_merge,
			delete_branch_on_merge=excluded.delete_branch_on_merge,
			allow_merge_commit=excluded.allow_merge_commit,
			subscribers_count=excluded.subscribers_count,
			network_count=excluded.network_count,
			organization_id=excluded.organization_id`

func storeRepository(dbConn *pgx.Conn, repo *github.Repository) error {
	owner := repo.GetOwner()
	var orgId *int64
	if repo.Organization != nil {
		orgId = repo.Organization.ID
	}

	_, err := dbConn.Exec(
		context.Background(), storeRepositoryQuery,
		repo.GetID(), repo.GetName(), repo.GetFullName(), owner.GetID(), owner.GetType(), repo.GetPrivate(),
		repo.GetDescription(), repo.GetFork(), repo.GetURL(), repo.GetForksCount(),
		repo.GetStargazersCount(), repo.GetWatchersCount(), repo.GetSize(), repo.GetDefaultBranch(),
		repo.GetOpenIssuesCount(), repo.GetIsTemplate(), repo.Topics, repo.GetHasIssues(),
		repo.GetHasProjects(), repo.GetHasWiki(), repo.GetHasPages(), repo.GetHasDownloads(),
		repo.GetArchived(), repo.GetDisabled(), repo.GetVisibility(),
		ghTsToTimeOrNull(repo.PushedAt), ghTsToTimeOrNull(repo.CreatedAt), ghTsToTimeOrNull(repo.UpdatedAt),
		repo.GetAllowRebaseMerge(), repo.GetAllowSquashMerge(), repo.GetDeleteBranchOnMerge(),
		repo.GetAllowMergeCommit(), repo.GetSubscribersCount(), repo.GetNetworkCount(), orgId)
	return err
}
