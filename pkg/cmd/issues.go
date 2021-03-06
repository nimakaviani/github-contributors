package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(issueCmd)
}

var (
	issueCmd = &cobra.Command{
		Use:   "issues",
		Short: "Analyze issues",
		Long:  `Analyze all contributors and companies creating the latest issues`,
		Run: func(cmd *cobra.Command, args []string) {

			if token := os.Getenv("GH_EMAIL_TOKEN"); token == "" {
				err := errors.New("GH_EMAIL_TOKEN needs to be set")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			s := scraper.NewGithubScraper("https://api.github.com")
			issues := analyzer.NewActivity(
				s,
				analyzer.NewCharter(s),
				models.Issue,
				count,
			)
			if err := issues.Process(repo); err != nil {
				println(err.Error())
				os.Exit(1)
			}
			if err := issues.Write(expand, output); err != nil {
				println(err.Error())
				os.Exit(1)
			}
		},
	}
)
