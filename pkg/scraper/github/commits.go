package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeCommitsOptions struct {
	OrgRepoAndRepoId
	// If true, also scrape the commit's change statistics (adds/mods/deletions).
	ScrapeStats bool
	// If true, scrape change stats for all files changed in the commit.
	ScrapeFiles bool
}

func ScrapeCommits(dbConn *pgx.Conn, ghClient *github.Client, options ScrapeCommitsOptions) (int, error) {
	// The ListCommits() call does not return the repo's ID. We'll query it separately.
	if err := ensureRepoIdFromAPI(ghClient, &options.OrgRepoAndRepoId); err != nil {
		return 0, nil
	}

	listAllOpts := github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	totalCommits := 0
	for {
		orgs, response, err := ghClient.Repositories.ListCommits(
			context.Background(), options.Owner, options.Repo, &listAllOpts)
		if err != nil {
			return totalCommits, err
		}
		for _, repoCommit := range orgs {
			if err := storeCommit(dbConn, options.RepoId, repoCommit); err != nil {
				return totalCommits, err
			}
			if options.ScrapeStats {
				statsOptions := ScrapeCommitStatsOptions{
					Owner:       options.Owner,
					Repo:        options.Repo,
					RepoId:      options.RepoId,
					CommitSHA:   repoCommit.GetSHA(),
					ScrapeFiles: options.ScrapeFiles,
				}
				if err := ScrapeCommitStats(dbConn, ghClient, statsOptions); err != nil {
					return totalCommits, err
				}
			}
			totalCommits += 1
		}
		if !continuePaginating(response) {
			break
		}
		listAllOpts.Page = response.NextPage
	}
	return totalCommits, nil
}

const storeCommitQuery = `
	INSERT INTO github_commits (
		repo_id, sha, author_id, committer_id, parents_sha, 
		commit_author_name, commit_author_email, commit_author_date,
		commit_committer_name, commit_committer_email, commit_committer_date,
		message, tree_sha, verification_verified, verification_reason)
	VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
	) ON CONFLICT (repo_id, sha)
		DO UPDATE SET 
			repo_id=excluded.repo_id,
			sha=excluded.sha,
			author_id=excluded.author_id,
			committer_id=excluded.committer_id,
			parents_sha=excluded.parents_sha,
			commit_author_name=excluded.commit_author_name,
			commit_author_email=excluded.commit_author_email,
			commit_author_date=excluded.commit_author_date,
			commit_committer_name=excluded.commit_committer_name,
			commit_committer_email=excluded.commit_committer_email,
			commit_committer_date=excluded.commit_committer_date,
			message=excluded.message,
			tree_sha=excluded.tree_sha,
			verification_verified=excluded.verification_verified,
			verification_reason=excluded.verification_reason`

func storeCommit(dbConn *pgx.Conn, repoId int64, repoCommit *github.RepositoryCommit) error {
	gitCommit := repoCommit.GetCommit()
	gitCommitAuthor := gitCommit.GetAuthor()
	gitCommitCommitter := gitCommit.GetCommitter()
	gitCommitVerification := gitCommit.GetVerification()
	parentsSha := make([]string, 0)
	var committerId *int64
	if repoCommit.Committer != nil {
		committerId = repoCommit.Committer.ID
	}
	for _, parent := range repoCommit.Parents {
		parentsSha = append(parentsSha, *parent.SHA)
	}

	_, err := dbConn.Exec(
		context.Background(), storeCommitQuery,
		repoId, repoCommit.GetSHA(), repoCommit.GetAuthor().GetID(), committerId, parentsSha,
		gitCommitAuthor.GetName(), gitCommitAuthor.GetEmail(), gitCommitAuthor.Date,
		gitCommitCommitter.GetName(), gitCommitCommitter.GetEmail(), gitCommitCommitter.Date,
		gitCommit.GetMessage(), gitCommit.GetTree().GetSHA(), gitCommitVerification.GetVerified(),
		gitCommitVerification.GetReason())
	return err
}
