/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/Kei-K23/git-go/internal/utils"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commits log",
	Run: func(cmd *cobra.Command, args []string) {
		// Get latest commit
		latestCommit := utils.GetCurrentCommit()

		// E.g Output for -log command
		// commit 20b94005ce5cc553525322d91b8eb3c9b7c79532 (HEAD -> main, origin/main)
		// Author: Author-A22 <aaa2301@gmail.com>
		// Date:   Tue Sep 24 15:54:09 2024 +0630

		// Format commit content to use in -log command
		// Track for parent tree
		currentCommit := latestCommit

		for currentCommit != "" {
			var parentCommit string
			// Get the current commit object
			currentCommitBytes := utils.ReadCommitObject(currentCommit)

			parts := strings.Split(string(currentCommitBytes), "\n")

			for _, part := range parts {

			}
		}

		// Recursively call commit object by following parent commit tress

		// Show human-readable commit log messages
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
