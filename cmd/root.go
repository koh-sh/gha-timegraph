/*
Copyright © 2023 koh-sh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/koh-sh/gha-timegraph/internal/gha"
	"github.com/koh-sh/gha-timegraph/internal/plotpng"
	"github.com/koh-sh/gha-timegraph/internal/stdout"
	"github.com/spf13/cobra"
)

// options
var (
	owner    string
	repo     string
	workflow string
	branch   string
	out      string
	outfile  string
	count    int
	silent   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gha-timegraph",
	Short: "Graphs the execution time of GitHub Actions",
	Long: `Graphs the execution time of GitHub Actions.

It creates GitHub Actions execution time as PNG graph.
Set GITHUB_TOKEN for private repositories.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Construct a new GitHub client
		client := gha.RtnClient()

		// parameters
		status := "success"
		if count < 0 {
			log.Fatal("count shall be bigger than 0")
		}

		runs, err := gha.GetRuns(client, count, owner, repo, workflow, branch, status, silent)
		if err != nil {
			log.Fatal(err)
		}
		if out == "png" {
			err := plotpng.SavePng(runs, outfile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("PNG saved to %s\n", outfile)
		} else {
			if e := stdout.PrintRuns(runs, out); e != nil {
				log.Fatal(e)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gha-timegraph.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVar(&owner, "owner", "", "Owner of the Action (Required)")
	rootCmd.Flags().StringVar(&repo, "repo", "", "Repository of the Action (Required)")
	rootCmd.Flags().StringVar(&workflow, "workflow", "", "workflow filename of the Action (Required)")
	rootCmd.Flags().StringVar(&branch, "branch", "", "Branch name to filter results")
	rootCmd.Flags().StringVar(&out, "out", "png", "format of output (csv or png)")
	rootCmd.Flags().StringVar(&outfile, "outfile", "graph.png", "name of output png file")
	rootCmd.Flags().IntVar(&count, "count", 30, "count of Workflow runs")
	rootCmd.Flags().BoolVar(&silent, "silent", false, "Hide Progress bar for GitHub API")
	rootCmd.MarkFlagRequired("owner")
	rootCmd.MarkFlagRequired("repo")
	rootCmd.MarkFlagRequired("workflow")
}

// set version from goreleaser variables
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}
