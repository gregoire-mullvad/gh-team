/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/google/go-github/v53/github"
	"github.com/spf13/cobra"
)

// syncLabelsCmd represents the syncLabels command
var syncLabelsCmd = &cobra.Command{
	Use:   "sync-labels",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		org, slug, err := parseTeam(team)
		if err != nil {
			return err
		}
		fmt.Println("syncLabels called")
		client, err := newClient()
		if err != nil {
			return err
		}
		ghRepos, err := listRepos(client, ctx, org, slug)

		term := term.FromEnv()
		termWidth, _, _ := term.Size()

		repos := make([]repository.Repository, 0, len(ghRepos))
		for _, r := range ghRepos {
			repos = append(repos, repository.Repository{Owner: r.GetOwner().GetLogin(), Name: r.GetName()})
		}

		// Collect existing labels
		state := make(map[repository.Repository]map[string]LabelSpec)
		for _, repo := range repos {
			labels, _, err := client.Issues.ListLabels(ctx, repo.Owner, repo.Name, nil)
			if err != nil {
				return err
			}
			state[repo] = make(map[string]LabelSpec, len(labels))
			for _, l := range labels {
				state[repo][l.GetName()] = LabelSpec{
					Name:        l.GetName(),
					Color:       l.GetColor(),
					Description: l.GetDescription(),
				}
			}
		}

		// Merge labels
		specs := make(map[string]LabelSpec)
		for _, labels := range state {
			for _, l := range labels {
				spec, ok := specs[l.Name]
				if !ok {
					spec.Name = l.Name
					spec.Color = l.Color
				}
				if spec.Description == "" && l.Description != "" {
					spec.Description = l.Description
				}
				specs[l.Name] = spec
			}
		}
		//				// color := label.GetColor()
		//				// table.AddField(label.GetName())
		//				// table.AddField(fmt.Sprintf("#%s", color))
		//				// table.AddField("█", tableprinter.WithColor(hexColor(color)))
		//				// table.AddField(label.GetDescription())
		//				// table.EndRow()
		//				specs[spec.Name] = spec
		//			}
		//			table.Render()
		//			fmt.Println()
		//		}

		// Determine the list of changes
		var changes []change
		for repo, repoLabels := range state {
			for name, label := range specs {
				repoLabel, ok := repoLabels[name]
				if !ok {
					changes = append(changes, change{repo: repo, new: label})
					continue
				}
				if reflect.DeepEqual(label, repoLabel) {
					continue
				}
				changes = append(changes, change{
					repo: repo,
					old:  repoLabel,
					new:  repoLabel,
				})
			}
		}

		table := tableprinter.New(term.Out(), term.IsTerminalOutput(), termWidth)
		for _, c := range changes {
			table.AddField(c.repo.Name)
			if c.old.Name == "" {
				table.AddField("CREATE")
			} else {
				table.AddField("UPDATE")
			}
			table.AddField(c.new.Name)
			if c.old.Color != c.new.Color {
				table.AddField(fmt.Sprintf("#%s", c.new.Color))
				table.AddField("█", tableprinter.WithColor(hexColor(c.new.Color)))
			} else {
				table.AddField("")
				table.AddField("")
			}
			if c.old.Description != c.new.Description {
				table.AddField(c.new.Description)
			} else {
				table.AddField("")
			}
			table.EndRow()
		}
		table.Render()
		fmt.Println()
		return nil
	},
}

type change struct {
	repo     repository.Repository
	old, new LabelSpec
}

func hexColor(h string) func(string) string {
	b, _ := hex.DecodeString(h)
	if len(b) != 3 {
		return func(s string) string { return s }
	}
	return func(s string) string { return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", b[0], b[1], b[2], s) }
}

type LabelSpec struct {
	Name        string
	Color       string
	Description string
}

func collectLabels(ctx context.Context, client *github.Client, repos []*github.Repository) error {
	return nil
}

func init() {
	rootCmd.AddCommand(syncLabelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncLabelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncLabelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
