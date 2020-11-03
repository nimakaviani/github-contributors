package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/spf13/cobra"
)

var (
	repo   string
	expand bool

	rootCmd = &cobra.Command{
		Use:   "github-contrib",
		Short: "github-contrib fetches contribution info for github",

		Run: func(cmd *cobra.Command, args []string) {

			if token := os.Getenv("GH_EMAIL_TOKEN"); token == "" {
				err := errors.New("GH_EMAIL_TOKEN needs to be set")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			scraper.Log("> pulling data from repo", repo)
			users, err := scraper.GetContribs(repo)
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}

			charter := analyzer.NewCharter()

			scraper.Log("> building charter ...")
			for _, user := range users {
				scraper.Log("   user", user.Login)
				err := charter.Build(user.Login)
				if err != nil {
					println(err.Error())
				}
			}

			scraper.Log("> done")
			scraper.Log(">> RESULTS")
			charter.Write(expand)
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
}
