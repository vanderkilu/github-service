package service

import (
	"context"
	"time"

	"github.com/google/go-github/v63/github"
	"github.com/vanderkilu/github-service/dao/postgresql"
)

type Github interface {
	CreateRepository(ctx context.Context, owner, repo string) error
	ProcessCommits(ctx context.Context, owner, repo string, since *time.Time) error
	MonitorRepo(ctx context.Context, owner, repo string, since *time.Time) error
}



type service struct {
	querier postgresql.Querier
	githubClient *github.Client
}

// New return service struct
func New(querier postgresql.Querier, githubClient *github.Client) *service {
	return &service{
		querier:        querier,
		githubClient:  githubClient,
	}
}