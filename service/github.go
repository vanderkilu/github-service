package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v63/github"
	"github.com/vanderkilu/github-service/dao/postgresql"
	"gopkg.in/guregu/null.v4"
)

type RepoCommit struct {
	SHA string
	Message string
	Author string
	URL string
	Date time.Time
}

func parseCommit(commit *github.RepositoryCommit) RepoCommit {
	repoCommit := RepoCommit{SHA: *commit.SHA, URL: *commit.URL}
	if (commit.Commit != nil) {
		repoCommit.Message = *commit.Commit.Message
		repoCommit.Author = *commit.Commit.Author.Email
		repoCommit.Date = commit.Commit.Author.Date.Time
	}
	return repoCommit
} 

func (svc *service) CreateRepository(ctx context.Context, owner, repo string) error {
	repository, _,  err := svc.githubClient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("GetRepositoryFromGithubError: %w", err)
	}
	err = svc.querier.CreateRepository(ctx, postgresql.CreateRepositoryParams{
		Description: *repository.Description,
		Url: *repository.URL,
		Language: *repository.Language,
		RepoName: *repository.Name,
		RepoFullName: *repository.FullName,
		ForksCount:  null.IntFrom(int64(*repository.ForksCount)),
		StarsCount: null.IntFrom(int64(*repository.StargazersCount)),
		OpenIssuesCount: null.IntFrom(int64(*repository.OpenIssuesCount)),
		WatchersCount: null.IntFrom(int64(*repository.WatchersCount)),
		CreatedAt: repository.CreatedAt.Time,
		UpdatedAt: repository.UpdatedAt.Time,
	})

	if err != nil {
		return fmt.Errorf("CreateRepository: %w", err)
	}
	fmt.Println("Repository saved successfully")
	return nil
}

func (svc *service) ProcessCommits(ctx context.Context, owner, repo string, since *time.Time) error {
	opt := github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// Get last commit, if it is the first time 
	// application is running sha will be empty
	// we will use this last commit sha to avoid pulling the 
	// same commit twice 
	sha, err := svc.querier.GetLastCommitSha(ctx)
	if err != nil {
		return fmt.Errorf("GetLastCommitShaError: %w", err)
	}

	//configurable param to enable us pull commits from that sha
	if sha != "" {
		opt.SHA = sha
		fmt.Println("Starting from SHA...", sha)
	}

	//configurable time since to pull commits
	if since != nil {
		opt.Since = *since
	}
	//handle pagination 
	for {
		commits, resp, err := svc.githubClient.Repositories.ListCommits(ctx, owner, repo, &opt)
		//handle rate limits
		if err != nil {
			if rateLimitError, ok := err.(*github.RateLimitError); ok {
				// Calculate the sleep duration until the rate limit resets
				now := time.Now()
				resetTime := rateLimitError.Rate.Reset
				sleepDuration := resetTime.Sub(now)          

				log.Printf("Rate limit exceeded. Sleeping for %v...\n", sleepDuration)
				time.Sleep(sleepDuration)
				continue
			}
			return fmt.Errorf("ListRepositoryCommitsError: %w", err)
		}

		//exit when there is no more commits to read
		if resp.NextPage == 0 {
			break
		}

		for _, commit := range commits {

			//This stores only unique commits(using the sha)
			repoCommit := parseCommit(commit)
			err= svc.querier.CreateCommit(ctx, postgresql.CreateCommitParams{
				Sha: repoCommit.SHA,
				RepoFullName: fmt.Sprintf("%s/%s", owner, repo),
				Message: repoCommit.Message,
				Author: repoCommit.Author,
				Url: repoCommit.URL,
				Date: repoCommit.Date,
			})
			if err != nil {
				return fmt.Errorf("CreateCommitError: %w", err)
			}
			fmt.Printf("Processed commit %s - %s successfully\n", *commit.SHA, *commit.Commit.Message)
		}
		opt.Page = resp.NextPage
	}
	return nil
}

func(svc *service) MonitorRepo(ctx context.Context, owner, repo string, since *time.Time) error {
	err := svc.CreateRepository(ctx, owner, repo)
	if err != nil {
		return err
	}
	return svc.ProcessCommits(ctx, owner, repo, since)
}