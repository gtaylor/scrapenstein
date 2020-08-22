package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapeIssuesOptions struct {
	OrgRepoAndRepoId
	// If true, scrape the PR summary stats (comment count, additions, deletions, etc).
	// This is VERY slow since we end up making one API request per scraped PR.
	ScrapeStats bool
}

func ScrapeIssues(dbConn *pgx.Conn, ghClient *github.Client, options ScrapeIssuesOptions) (int, error) {
	// The ListByRepo() call does not return the repo's ID. We'll query it separately.
	if err := ensureRepoIdFromAPI(ghClient, &options.OrgRepoAndRepoId); err != nil {
		return 0, nil
	}

	listOpts := github.IssueListByRepoOptions{State: "all"}
	totalIssues := 0
	for {
		issues, response, err := ghClient.Issues.ListByRepo(context.Background(), options.Owner, options.Repo, &listOpts)
		if err != nil {
			return totalIssues, err
		}
		for _, issue := range issues {
			if issue.IsPullRequest() {
				continue
			}
			if err := storeIssue(dbConn, options.RepoId, issue); err != nil {
				return totalIssues, err
			}
			totalIssues += 1
		}
		if !continuePaginating(response) {
			break
		}
		listOpts.Page = response.NextPage
	}
	return totalIssues, nil
}

const storeIssueQuery = `
	INSERT INTO github_issues (
		id, repo_id, number, state, locked, title, user_id, created_at, updated_at,
		closed_at, labels, comments, assignee_ids)
	VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	) ON CONFLICT (id)
		DO UPDATE SET 
			id=excluded.id, 
			repo_id=excluded.repo_id, 
			number=excluded.number,
			state=excluded.state,
			locked=excluded.locked,
			title=excluded.title,
			user_id=excluded.user_id,
			created_at=excluded.created_at,
			updated_at=excluded.updated_at,
			closed_at=excluded.closed_at,
			labels=excluded.labels,
			comments=excluded.comments,
			assignee_ids=excluded.assignee_ids
`

func storeIssue(dbConn *pgx.Conn, repoId int64, issue *github.Issue) error {
	labels := make([]string, 0)
	for _, label := range issue.Labels {
		labels = append(labels, label.GetName())
	}
	assigneeIds := make([]int64, 0)
	for _, assignee := range issue.Assignees {
		assigneeIds = append(assigneeIds, *assignee.ID)
	}

	_, err := dbConn.Exec(
		context.Background(), storeIssueQuery,
		issue.GetID(), repoId, issue.GetNumber(), issue.GetState(), issue.GetLocked(), issue.GetTitle(),
		issue.GetUser().GetID(), issue.GetCreatedAt(), issue.GetUpdatedAt(), issue.ClosedAt,
		labels, issue.GetComments(), assigneeIds)
	return err
}
