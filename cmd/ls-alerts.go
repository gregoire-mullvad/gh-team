/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/google/go-github/v53/github"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

var (
	lsAlertsState     string
	lsAlertsScope     string
	lsAlertsEcosystem string
)

// lsAlertsCmd represents the lsAlerts command
var lsAlertsCmd = &cobra.Command{
	Use:   "ls-alerts",
	Short: "List dependabot alerts in the team's repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		org, slug, err := parseTeam(team)
		if err != nil {
			return err
		}
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Color("reset")
		s.Suffix = fmt.Sprintf(" Loading alerts for @%s", team)
		s.Start() // Start the spinner
		defer s.Stop()
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
		errTable := tableprinter.New(term.Out(), term.IsTerminalOutput(), termWidth)
		opts := github.ListAlertsOptions{}
		if lsAlertsState != "" {
			opts.State = &lsAlertsState
		}
		if lsAlertsScope != "" {
			opts.Scope = &lsAlertsScope
		}
		if lsAlertsEcosystem != "" {
			opts.Ecosystem = &lsAlertsEcosystem
		}
		var alerts []*github.DependabotAlert
		for _, repo := range repos {
			repoAlerts, _, err := client.Dependabot.ListRepoAlerts(
				context.TODO(),
				repo.GetOwner().GetLogin(),
				repo.GetName(),
				&opts,
			)
			if err != nil {
				errTable.AddField(repo.GetFullName(), tableprinter.WithColor(ansi.ColorFunc("red+b")))
				var errr *github.ErrorResponse
				if errors.As(err, &errr) {
					errTable.AddField(errr.Message, tableprinter.WithColor(ansi.ColorFunc("red+b")))
				} else {
					errTable.AddField(err.Error(), tableprinter.WithColor(ansi.ColorFunc("red+b")))
				}
				errTable.EndRow()
				continue
			}
			for _, a := range repoAlerts {
				a.Repository = repo
				alerts = append(alerts, a)
			}
		}
		s.Stop()

		for _, alert := range alerts {
			table.AddField(
				strings.ToUpper(alert.GetSecurityVulnerability().GetSeverity()),
				tableprinter.WithColor(ansi.ColorFunc(colorForAlertSeverity(alert))),
			)
			table.AddField(alert.GetRepository().GetFullName(), tableprinter.WithColor(ansi.ColorFunc("gray+b")))
			table.AddField(
				fmt.Sprintf("#%d", *alert.Number),
			)
			table.AddField(alert.GetDependency().GetPackage().GetName())
			table.AddField(fmt.Sprintf("%s (%s)",
				alert.GetDependency().GetManifestPath(),
				alert.GetDependency().GetPackage().GetEcosystem(),
			))
			table.AddField(alert.GetDependency().GetScope())
			table.AddField(
				fmt.Sprintf("%s, %s", strings.Title(alert.GetState()),
					text.RelativeTimeAgo(time.Now(), (alert.GetCreatedAt()).Time),
				),
				tableprinter.WithColor(ansi.ColorFunc("gray")))
			table.EndRow()
		}
		if err := table.Render(); err != nil {
			return err
		}
		if err := errTable.Render(); err != nil {
			return err
		}
		return nil
	},
}

// colorForAlertSeverity returns a color that depends on the severity of the alert.
func colorForAlertSeverity(a *github.DependabotAlert) string {
	switch a.GetSecurityVulnerability().GetSeverity() {
	case "critical":
		return "red"
	case "high":
		return "red+h"
	case "medium":
		return "yellow"
	case "low":
		return "blue"
	default:
		return ""
	}
}

func init() {
	rootCmd.AddCommand(lsAlertsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsAlertsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	lsAlertsCmd.Flags().StringVar(&lsAlertsState, "state", "",
		"A comma-separated list of states. If specified, only alerts with these states will be returned (auto_dismissed, dismissed, fixed, open)")
	lsAlertsCmd.Flags().StringVar(&lsAlertsScope, "scope", "",
		"The scope of the vulnerable dependency. If specified, only alerts with this scope will be returned (development, runtime)")
	lsAlertsCmd.Flags().StringVar(&lsAlertsEcosystem, "ecosystem", "",
		"A comma-separated list of ecosystems. If specified, only alerts for these ecosystems will be returned (composer, go, maven, npm, nuget, pip, pub, rubygems, rust)")
}
