package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

var authorFilter string

// lsPullsCmd represents the lsPulls command
var lsPullsCmd = &cobra.Command{
	Use:     "ls-prs",
	Aliases: []string{"ls-pulls"},
	Short:   "List open pull requests in the team's repositories",
	Long: `List open pull requests in the team's repositories.  Example:
    gh team ls-prs
    myorg/repo1  #123  Add the new thing                   mybranch  someuser, about 1 week ago
    myorg/repo2  #42   Life, the universe, and everything  answer    deepthought, about 1 million years ago

It will only print PRs from repos the team can push to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		org, slug, err := parseTeam(team)
		if err != nil {
			return err
		}
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Color("reset")
		s.Suffix = fmt.Sprintf(" Loading PRs for @%s", team)
		s.Start() // Start the spinner
		client, err := newClient()
		if err != nil {
			return err
		}
		repos, err := listRepos(client, context.TODO(), org, slug)
		if err != nil {
			return err
		}
		term := term.FromEnv()
		termWidth, _, _ := term.Size()
		table := tableprinter.New(term.Out(), term.IsTerminalOutput(), termWidth)
		for _, repo := range repos {
			pulls, err := listOpenPullRequests(client, context.TODO(), *repo.Owner.Login, *repo.Name)
			if err != nil {
				return err
			}

			for _, pull := range pulls {
				if authorFilter != "" && *pull.User.Login != authorFilter {
					continue
				}
				table.AddField(*repo.FullName, tableprinter.WithColor(ansi.ColorFunc("gray+b")))
				table.AddField(
					fmt.Sprintf("#%d", *pull.Number),
					tableprinter.WithColor(ansi.ColorFunc(colorForPRState(pull))),
				)
				table.AddField(*pull.Title)
				table.AddField(*pull.Head.Ref, tableprinter.WithColor(ansi.ColorFunc("cyan")))
				table.AddField(
					fmt.Sprintf("%s, %s", *pull.User.Login,
						text.RelativeTimeAgo(time.Now(), (*pull.CreatedAt).Time),
					),
					tableprinter.WithColor(ansi.ColorFunc("gray")))
				table.EndRow()
			}
		}
		s.Stop()
		if err := table.Render(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	lsPullsCmd.Flags().StringVarP(&authorFilter, "author", "A", "", "Filter by author")
	rootCmd.AddCommand(lsPullsCmd)
}
