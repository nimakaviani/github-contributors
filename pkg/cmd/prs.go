package cmd

import (
	"os"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/models"
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
			charter := analyzer.NewCharter()
			issues := analyzer.NewActivity(models.PR, charter, count)
			if err := issues.Process(repo); err != nil {
				println(err.Error())
				os.Exit(1)
			}
			issues.Write()
		},
	}
)
