package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var lsReposShowRemotes bool

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
		for _, repo := range repos {
			if lsReposShowRemotes {
				fmt.Println(*repo.SSHURL)
			} else {
				fmt.Println(*repo.FullName)
			}
		}
		return nil
	},
}

func init() {
	lsReposCmd.Flags().BoolVarP(&lsReposShowRemotes, "remotes", "r", false, "Print git remotes instead of repository names")
	rootCmd.AddCommand(lsReposCmd)
}
