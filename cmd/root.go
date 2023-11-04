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
	"github.com/spf13/cobra"
)

// options
var (
	owner    string
	repo     string
	filename string
	branch   string
	out      string
	count    int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gha-timegraph",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Construct a new GitHub client
		client := gha.RtnClient()

		// parameters
		status := "success"

		runs, err := gha.GetRuns(client, count, owner, repo, filename, branch, status)
		if err != nil {
			log.Fatal(err)
		}
		if out == "csv" {
			fmt.Printf("%s,%s,%s\n", "Name", "StartTime(UTC)", "Elapsed")
			for _, v := range runs {
				fmt.Printf("%s,%s,%g\n", v.Name, v.Starttime.Format("2006-01-02 15:04:05"), v.Elapsed)
			}
		} else if out == "png" {
			err := plotpng.SavePng(runs)
			if err != nil {
				log.Fatal(err)
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
	rootCmd.Flags().StringVar(&filename, "filename", "", "Filename of the Action (Required)")
	rootCmd.Flags().StringVar(&branch, "branch", "", "Branch name to filter results")
	rootCmd.Flags().StringVar(&out, "out", "csv", "format of output (csv or png)")
	rootCmd.Flags().IntVar(&count, "count", 30, "count of Workflow runs")
	rootCmd.MarkFlagRequired("owner")
	rootCmd.MarkFlagRequired("repo")
	rootCmd.MarkFlagRequired("filename")
}

// set version from goreleaser variables
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}
