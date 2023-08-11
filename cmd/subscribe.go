/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

// subscribe.goCmd represents the subscribe.go command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to the team repostories",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		subscriptions, err := listSubscriptions(client, context.TODO())
		if err != nil {
			return err
		}
		if subscribeCmdCheck {
			term := term.FromEnv()
			termWidth, _, _ := term.Size()
			table := tableprinter.New(term.Out(), term.IsTerminalOutput(), termWidth)
			for _, repo := range repos {
				subscribed := subscriptions[repo.GetID()]
				table.AddField(repo.GetFullName(), tableprinter.WithColor(ansi.ColorFunc("gray+b")))
				if subscribed {
					table.AddField("👁️", tableprinter.WithColor(ansi.ColorFunc("green")))
				} else {
					table.AddField("🚫", tableprinter.WithColor(ansi.ColorFunc("red")))
				}
				table.EndRow()
			}
			if err := table.Render(); err != nil {
				return err
			}
		} else {
			for _, repo := range repos {
				subscribed := subscriptions[repo.GetID()]
				if !subscribed {
					fmt.Printf("Subscribing to %s\n", repo.GetFullName())
					err := subscribe(client, context.TODO(), *repo.Owner.Login, repo.GetName())
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	},
}

var subscribeCmdCheck bool

func init() {
	rootCmd.AddCommand(subscribeCmd)
	subscribeCmd.Flags().BoolVarP(&subscribeCmdCheck, "check", "c", false, "Check but don't modify subscriptions")
}
