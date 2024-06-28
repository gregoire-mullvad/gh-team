package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-team [--team <org>/<team>] [command]",
	Short: "Work with Github teams.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var team string
var minPermission string
var excludeReposRegexp string

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&team, "team", "T", os.ExpandEnv("$GITHUB_TEAM"), "github team",
	)
	rootCmd.PersistentFlags().StringVarP(
		&minPermission, "min-permissions", "P", "push", "only include repos for which the team has at least this permission",
	)
	rootCmd.PersistentFlags().StringVarP(
		&excludeReposRegexp, "exclude", "E", os.ExpandEnv("$GITHUB_TEAM_EXCLUDE"),
		"exclude repositories that match the given regular expression",
	)
}
