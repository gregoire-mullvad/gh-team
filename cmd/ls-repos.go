package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/config"
	"github.com/cli/go-gh/v2/pkg/template"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/spf13/cobra"
)

var lsReposShowRemotes bool
var lsReposLanguage string
var lsReposTemplate string

// lsReposCmd represents the lsRepos command
var lsReposCmd = &cobra.Command{
	Use:   "ls-repos",
	Args:  cobra.NoArgs,
	Short: "List the team's repositories",
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
		if lsReposTemplate != "" {
			term := term.FromEnv()
			tWidth, _, _ := term.Size()
			tmpl := template.New(os.Stdout, tWidth, term.IsColorEnabled())
			if err := tmpl.Parse(lsReposTemplate); err != nil {
				return err
			}
			json, err := json.Marshal(repos)
			if err != nil {
				return err
			}
			return tmpl.Execute(bytes.NewReader(json))
		}
		proto := getGitProtocol()
		for _, repo := range repos {
			if lsReposLanguage != "" && !strings.EqualFold(repo.GetLanguage(), lsReposLanguage) {
				continue
			}
			if lsReposTemplate != "" {
			} else if lsReposShowRemotes && proto == "https" {
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
	lsReposCmd.Flags().StringVarP(&lsReposTemplate, "template", "t", "", `Format JSON output using a Go template; see "gh help formatting"`)
	rootCmd.AddCommand(lsReposCmd)
}
