package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/cli/go-gh/pkg/config"
	"github.com/spf13/cobra"
)

var lsReposShowRemotes bool
var lsReposLanguage string

// lsReposCmd represents the lsRepos command
var lsReposCmd = &cobra.Command{
	Use:   "ls-repos",
	Args:  cobra.NoArgs,
	Short: "Print the team's repositories to stdout",
	Long: `Print the team's repositories to stdout. Example:

    gh team ls-repos myorg/myteam
    myorg/repo1
    myorg/repo2

It will only print repos the team can push to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		org, slug, err := parseTeam(team)
		if err != nil {
			return err
		}
		client, err := newClient()
		if err != nil {
			return err
		}
		repos, err := listRepos(client, context.TODO(), org, slug)
		if err != nil {
			return err
		}
		proto := getGitProtocol()
		for _, repo := range repos {
			if lsReposLanguage != "" && !strings.EqualFold(repo.GetLanguage(), lsReposLanguage) {
				continue
			}
			if lsReposShowRemotes && proto == "https" {
				fmt.Println(*repo.CloneURL)
			} else if lsReposShowRemotes {
				fmt.Println(*repo.SSHURL)
			} else {
				fmt.Println(*repo.FullName)
			}
		}
		return nil
	},
}

func getGitProtocol() string {
	config, err := config.Read()
	if err != nil {
		return "ssh"
	}
	if proto, err := config.Get([]string{"hosts", "github.com", "git_protocol"}); err == nil {
		return proto
	}
	if proto, err := config.Get([]string{"git_protocol"}); err == nil {
		return proto
	}
	return "ssh"
}

func init() {
	lsReposCmd.Flags().StringVarP(&lsReposLanguage, "language", "l", "", "Filter by primary coding language")
	lsReposCmd.Flags().BoolVarP(&lsReposShowRemotes, "remotes", "r", false, "Print git remotes instead of repository names")
	rootCmd.AddCommand(lsReposCmd)
}
