package cmd

import (
	"os"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(prCmd)
}

var (
	prCmd = &cobra.Command{
		Use:   "prs",
		Short: "Analyze PRs",
		Long:  `Analyze all contributors and companies submitting pull requests to the repo`,
		Run: func(cmd *cobra.Command, args []string) {
			s := scraper.NewGithubScraper("https://api.github.com")
			prs := analyzer.NewActivity(
				s,
				analyzer.NewCharter(s),
				models.PR,
				count,
			)
			if err := prs.Process(repo); err != nil {
				println(err.Error())
				os.Exit(1)
			}
			if err := prs.Write(expand, output); err != nil {
				println(err.Error())
				os.Exit(1)
			}
		},
	}
)
