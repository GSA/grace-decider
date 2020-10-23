package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func newAuthenticatedClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func (r *req) getPR() (*github.PullRequest, error) {
	ctx := context.Background()

	pr, _, err := r.githubClient.PullRequests.Get(ctx, r.owner, r.repoName, r.prNumber)
	if err != nil {
		return pr, err
	}

	return pr, nil
}

func waitForMerge(pr *github.PullRequest) error {
	client := newAuthenticatedClient()
	ctx := context.Background()
	owner := *pr.Base.Repo.Owner.Login
	repo := *pr.Base.Repo.Name
	var err error

	fmt.Print("Waiting for Pull Request to be merged")
	for *pr.State != "closed" {
		fmt.Print(".")
		time.Sleep(10 * time.Second)
		pr, _, err = client.PullRequests.Get(ctx, owner, repo, *pr.Number)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	if !*pr.Merged {
		return fmt.Errorf("pull request %s but not merged", *pr.State)
	}

	return nil
}
