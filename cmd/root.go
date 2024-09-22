/*
Copyright © 2024 Kei-K23 <arkar.dev.kei@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-go",
	Short: "A Git implementation in Go",
	Long:  `git-go is a lightweight version of Git implemented in Go.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to git-go! Use 'git-go --help' for usage information.")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
