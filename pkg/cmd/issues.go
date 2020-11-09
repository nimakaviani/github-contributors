package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.PersistentFlags().IntVarP(&count, "count", "c", 10, "count of issues to analyze")
}

var (
	count int

	issueCmd = &cobra.Command{
		Use:   "issues",
		Short: "Analyze issues",
		Long:  `Analyze all contributors and companies crreating the latest issues`,
		Run: func(cmd *cobra.Command, args []string) {

			if token := os.Getenv("GH_EMAIL_TOKEN"); token == "" {
				err := errors.New("GH_EMAIL_TOKEN needs to be set")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			bar := pb.StartNew(count)
			charter := analyzer.NewCharter()
			issues := analyzer.NewIssues(charter, count)
			if err := issues.Process(repo, bar); err != nil {
				println(err.Error())
				os.Exit(1)
			}
			bar.Finish()
			issues.Write()
		},
	}
)
