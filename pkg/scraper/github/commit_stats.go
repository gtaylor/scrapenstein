package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeCommitStatsOptions struct {
	Owner string
	Repo  string
	// Providing RepoID allows for skipping a query for the repo's ID.
	// If this value is 0, the repo's ID will be queried from the Owner+Repo combo.
	RepoId    int64
	CommitSHA string
	// If true, scrape change stats for all files changed in the commit.
	ScrapeFiles bool
}

func ScrapeCommitStats(dbConn *pgx.Conn, ghClient *github.Client, options ScrapeCommitStatsOptions) error {
	if options.RepoId == 0 {
		// The Repo's ID is used as our repo PKey instead of the owner + name
		// since repo names can change.
		repo, _, err := ghClient.Repositories.Get(context.Background(), options.Owner, options.Repo)
		if err != nil {
			return err
		}
		options.RepoId = repo.GetID()
	}

	repoCommit, _, err := ghClient.Repositories.GetCommit(
		context.Background(), options.Owner, options.Repo, options.CommitSHA)
	if err != nil {
		return err
	}
	if err := storeCommitStats(dbConn, options, repoCommit); err != nil {
		return err
	}
	return nil
}

const storeCommitStatsQuery = `
	INSERT INTO github_commit_stats (
		repo_id, sha, additions, deletions, total)
	VALUES(
		$1, $2, $3, $4, $5
	) ON CONFLICT (repo_id, sha)
		DO UPDATE SET 
			repo_id=excluded.repo_id,
			sha=excluded.sha,
			additions=excluded.additions,
			deletions=excluded.deletions,
			total=excluded.total`

const storeCommitFilesQuery = `
	INSERT INTO github_commit_files (
		repo_id, sha, filename, additions, deletions, changes, status)
	VALUES(
		$1, $2, $3, $4, $5, $6, $7
	) ON CONFLICT (repo_id, sha, filename)
		DO UPDATE SET 
			repo_id=excluded.repo_id,
			sha=excluded.sha,
			filename=excluded.filename,
			additions=excluded.additions,
			deletions=excluded.deletions,
			changes=excluded.changes,
			status=excluded.status`

func storeCommitStats(dbConn *pgx.Conn, options ScrapeCommitStatsOptions, repoCommit *github.RepositoryCommit) error {
	repoCommitStats := repoCommit.GetStats()

	_, err := dbConn.Exec(
		context.Background(), storeCommitStatsQuery,
		options.RepoId, repoCommit.GetSHA(), repoCommitStats.GetAdditions(), repoCommitStats.GetDeletions(),
		repoCommitStats.GetTotal())
	if err != nil {
		return err
	}

	// There are potentially a ton of these, so storing them is optional.
	if !options.ScrapeFiles {
		return nil
	}
	for _, commitFile := range repoCommit.Files {
		_, err = dbConn.Exec(
			context.Background(), storeCommitFilesQuery,
			options.RepoId, repoCommit.GetSHA(), commitFile.GetFilename(),
			commitFile.GetAdditions(), commitFile.GetDeletions(), commitFile.GetChanges(),
			commitFile.GetStatus())
	}
	return err
}
