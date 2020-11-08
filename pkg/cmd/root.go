package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
)

var (
	repo         string
	expand       bool
	enableIssues bool
	issueCount   int

	rootCmd = &cobra.Command{
		Use:   "github-contrib",
		Short: "github-contrib fetches contribution info for github",

		Run: func(cmd *cobra.Command, args []string) {

			if token := os.Getenv("GH_EMAIL_TOKEN"); token == "" {
				err := errors.New("GH_EMAIL_TOKEN needs to be set")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			charter := analyzer.NewCharter()
			// if err := charter.Process(repo); err != nil {
			// 	println(err.Error())
			// 	os.Exit(1)
			// }

			issues := analyzer.NewIssues(charter, enableIssues, issueCount)
			if err := issues.Process(repo); err != nil {
				println(err.Error())
				os.Exit(1)
			}

			charter.Write(expand)
			issues.Write()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "project repo")
	rootCmd.PersistentFlags().BoolVarP(&expand, "expand", "e", true, "expand user info")
	rootCmd.PersistentFlags().BoolVarP(&scraper.Anonymous, "unauthenticated", "u", false, "unauthenticated gh call")
	rootCmd.PersistentFlags().BoolVarP(&scraper.Debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().BoolVarP(&enableIssues, "issues", "i", false, "analyze issues")
	rootCmd.PersistentFlags().IntVarP(&issueCount, "count", "c", 10, "count of issues to analyze")
}
