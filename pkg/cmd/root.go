package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/nimakaviani/github-contributors/pkg/utils"
)

type Activity int

const (
	Issue = iota
	PR
)

var (
	repo   string
	output string
	expand bool
	count  int

	rootCmd = &cobra.Command{
		Use:   "github-contrib",
		Short: "fetch contribution info for a github repo.",

		Run: func(cmd *cobra.Command, args []string) {

			charter := analyzer.NewCharter(scraper.NewGithubScraper("https://api.github.com"))
			if err := charter.Process(repo, count); err != nil {
				println(err.Error())
				os.Exit(1)
			}

			if err := charter.Write(expand, output); err != nil {
				println(err.Error())
				os.Exit(1)
			}
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
	if token := os.Getenv("GH_EMAIL_TOKEN"); token == "" {
		msg := fmt.Sprintf(`Before you get started, "GH_EMAIL_TOKEN" needs to be set
The tool requires authenticated requests to retrieve contributor emails, see https://git.io/vxctz.
To get a token, visit https://github.com/settings/tokens/new?description=github-email.
You don't need to check any of the checkboxes. Just generate the token, and export it
in your terminal: "export GH_EMAIL_TOKEN=<token>"
		`)
		err := errors.New(msg)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "project repo")
	rootCmd.PersistentFlags().BoolVarP(&expand, "expand", "e", true, "expand user info")
	rootCmd.PersistentFlags().BoolVarP(&scraper.Anonymous, "unauthenticated", "u", false, "unauthenticated gh call")
	rootCmd.PersistentFlags().BoolVarP(&utils.Debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().IntVarP(&count, "count", "c", 30, "count of items to analyze")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format")
}
