package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/google/go-github/v53/github"
)

func parseTeam(team string) (string, string, error) {
	if team == "" {
		return "", "", errors.New("no team set (use --team or $GH_TEAM)")
	}
	org, slug, ok := strings.Cut(team, "/")
	if !ok {
		return "", "", fmt.Errorf("%s: invalid team, expected <org>/<team>", team)
	}
	return org, slug, nil
}

func newClient() (*github.Client, error) {
	httpclient, err := api.DefaultHTTPClient()
	if err != nil {
		return nil, err
	}
	return github.NewClient(httpclient), nil
}

func listRepos(client *github.Client, ctx context.Context, org, team string) ([]*github.Repository, error) {
	repos, _, err := client.Teams.ListTeamReposBySlug(ctx, org, team, nil)
	if err != nil {
		return nil, err
	}
	var result []*github.Repository
	for _, repo := range repos {
		if repo.Permissions["push"] {
			result = append(result, repo)
		}
	}
	return result, nil
}

func listOpenPullRequests(client *github.Client, ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	prs, _, err := client.PullRequests.List(ctx, owner, repo, nil)
	return prs, err
}

// colorForPRState returns a color that depends on the state of the PR.
// Use the same colors as gh.
func colorForPRState(pr *github.PullRequest) string {
	switch *pr.State {
	case "open":
		if *pr.Draft {
			return "gray"
		}
		return "green"
	case "closed":
		return "red"
	case "merged":
		return "magenta"
	default:
		return ""
	}
}
