package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/jackc/pgx/v4"
)

type ScrapePullRequestsOptions struct {
	OrgRepoAndRepoId
	// If true, scrape the PR summary stats (comment count, additions, deletions, etc).
	// This is VERY slow since we end up making one API request per scraped PR.
	ScrapeStats bool
}

func ScrapePullRequests(dbConn *pgx.Conn, ghClient *github.Client, options ScrapePullRequestsOptions) (int, error) {
	// The List() call does not return the repo's ID. We'll query it separately.
	if err := ensureRepoIdFromAPI(ghClient, &options.OrgRepoAndRepoId); err != nil {
		return 0, nil
	}

	listOpts := github.PullRequestListOptions{State: "all"}
	totalPRs := 0
	for {
		pullRequests, response, err := ghClient.PullRequests.List(context.Background(), options.Owner, options.Repo, &listOpts)
		if err != nil {
			return totalPRs, err
		}
		for _, pullRequest := range pullRequests {
			if options.ScrapeStats {
				// The API does not return comment, review comment, or change stats in List mode.
				// To get those stats, we must issue a Get request for each PR. :(
				detailedPR, _, err := ghClient.PullRequests.Get(context.Background(), options.Owner, options.Repo, pullRequest.GetNumber())
				if err != nil {
					return totalPRs, err
				}
				pullRequest = detailedPR
			}
			if err := storePullRequest(dbConn, options.RepoId, pullRequest); err != nil {
				return totalPRs, err
			}
			totalPRs += 1
		}
		if !continuePaginating(response) {
			break
		}
		listOpts.Page = response.NextPage
	}
	return totalPRs, nil
}

const storePullRequestQuery = `
	INSERT INTO github_pull_requests (
		id, repo_id, number, state, locked, title, user_id, created_at, updated_at,
		closed_at, merged_at, labels, draft, merged, mergeable, merged_by_id,
		rebaseable, comments, commits, additions, deletions, changed_files,
		review_comments, assignee_ids, requested_reviewer_ids)
	VALUES(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
		$19, $20, $21, $22, $23, $24, $25
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
			merged_at=excluded.merged_at,
			labels=excluded.labels,
			draft=excluded.draft,
			merged=excluded.merged,
			mergeable=excluded.mergeable,
			merged_by_id=excluded.merged_by_id,
			rebaseable=excluded.rebaseable,
			comments=excluded.comments,
			commits=excluded.commits,
			additions=excluded.additions,
			deletions=excluded.deletions,
			changed_files=excluded.changed_files,
			review_comments=excluded.review_comments,
			assignee_ids=excluded.assignee_ids,
			requested_reviewer_ids=excluded.requested_reviewer_ids
`

func storePullRequest(dbConn *pgx.Conn, repoId int64, pr *github.PullRequest) error {
	labels := make([]string, 0)
	for _, label := range pr.Labels {
		labels = append(labels, label.GetName())
	}
	var mergedBy *int64
	if pr.MergedBy != nil {
		mergedBy = pr.MergedBy.ID
	}
	assigneeIds := make([]int64, 0)
	for _, assignee := range pr.Assignees {
		assigneeIds = append(assigneeIds, *assignee.ID)
	}
	requestedReviewerIds := make([]int64, 0)
	for _, reviewer := range pr.RequestedReviewers {
		requestedReviewerIds = append(requestedReviewerIds, *reviewer.ID)
	}

	_, err := dbConn.Exec(
		context.Background(), storePullRequestQuery,
		pr.GetID(), repoId, pr.GetNumber(), pr.GetState(), pr.GetLocked(), pr.GetTitle(),
		pr.GetUser().GetID(), pr.GetCreatedAt(), pr.GetUpdatedAt(), pr.ClosedAt, pr.MergedAt,
		labels, pr.GetDraft(), pr.GetMerged(), pr.GetMergeable(), mergedBy, pr.GetRebaseable(),
		pr.GetComments(), pr.GetCommits(), pr.GetAdditions(), pr.GetDeletions(),
		pr.GetChangedFiles(), pr.GetReviewComments(), assigneeIds, requestedReviewerIds)
	return err
}
