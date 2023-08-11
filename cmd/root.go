package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-team [--team <org>/<team>] [command]",
	Short: "A brief description of your application",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var team string

func init() {
	rootCmd.PersistentFlags().StringVarP(&team, "team", "T", os.ExpandEnv("$GITHUB_TEAM"), "github team (default is $GITHUB_TEAM)")
}
