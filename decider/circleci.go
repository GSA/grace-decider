package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/jszwedko/go-circleci"
)

func newCircleClient(token string) *circleci.Client {
	return &circleci.Client{Token: token}
}

func waitForApply(pr *github.PullRequest) error {
	const sleepSec = 5
	fmt.Println("Waiting for CircleCI apply_terraform job to complete")
	client := &circleci.Client{Token: os.Getenv("CIRCLE_TOKEN")}
	account := *pr.Base.Repo.Owner.Login
	repo := *pr.Base.Repo.Name
	branch := *pr.Base.Ref
	sha := *pr.Head.SHA
	startTime := *pr.MergedAt // Only interested in builds that started after merge
	timeout := 5 * time.Minute
	const numJobs = -1 // 4 is Number of jobs in workflow
	var build *circleci.Build

	for build == nil {
		builds, err := client.ListRecentBuildsForProject(account, repo, branch, "", numJobs, 0)
		if err != nil {
			return err
		}

		for i, b := range builds {
			buildSHA := "unknown"
			if len(b.AllCommitDetails) > 0 {
				buildSHA = b.AllCommitDetails[0].Commit
			}
			fmt.Printf("%d) %d Job: %s SHA: %s Status: %s Lifecycle: %s Outcome: %s\n",
				i, b.BuildNum, b.BuildParameters["CIRCLE_JOB"], buildSHA,
				b.Status, b.Lifecycle, b.Outcome)
			if buildSHA == sha && b.StartTime.After(startTime) {
				if b.BuildParameters["CIRCLE_JOB"] == "apply_terraform" {
					timeout = 30 * time.Minute
					build = b
				}
				if build.Lifecycle != "finished" {
					err := waitForBuild(client, b, timeout)
					if err != nil {
						return err
					}
				}
				if *b.Failed {
					return fmt.Errorf("%s %s", b.BuildParameters["CIRCLE_JOB"], b.Status)
				}
			}
		}
		time.Sleep(sleepSec * time.Second)
	}

	return nil
}

// waitForBuild ... used internally to wait for the build matching the given
// buildNum to complete, does not validate that the build was successful
// jobTimeout is the duration to wait before giving up
func waitForBuild(client *circleci.Client, build *circleci.Build, jobTimeout time.Duration) (err error) {
	const sleepSec = 5
	var (
		count   int
		endTime = time.Now().Add(jobTimeout)
	)
	for {
		if time.Now().After(endTime) {
			return fmt.Errorf("job timeout exceeded while waiting for build %s [%d] to finish", build.BuildParameters["CIRCLE_JOB"], build.BuildNum)
		}
		if count%10 == 0 {
			fmt.Printf("waiting for build %s [%d] to finish\n", build.BuildParameters["CIRCLE_JOB"], build.BuildNum)
		}
		time.Sleep(sleepSec * time.Second)
		build, err = client.GetBuild(build.Username, build.Reponame, build.BuildNum)
		if err != nil {
			return err
		}
		// Lifecycle options:
		// :queued, :scheduled, :not_run, :not_running, :running or :finished
		if build.Lifecycle == "finished" {
			fmt.Printf("job %s [%d] %s with state %s\n", build.BuildParameters["CIRCLE_JOB"], build.BuildNum, build.Lifecycle, build.Status)
			return nil
		}
		count++
	}
}
