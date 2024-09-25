/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
			// Get the current commit object
			currentCommitBytes := utils.ReadCommitObject(currentCommit)
			currentCommitContent := string(currentCommitBytes)
			parts := strings.Split(currentCommitContent, "\n")

			var author, date, message, parentCommit string

			for _, part := range parts {
				// Parse author information
				if strings.HasPrefix(part, "author") && part != "" {
					// Example for author content (e.g author author <author@gmail.com> 2024-09-25T19:32:44+06:30)
					authorParts := strings.SplitN(part, " ", 3)
					author = authorParts[1] + " " + authorParts[2]
				}

				if strings.HasPrefix(part, "committer") && part != "" {
					// Example for committer content (e.g committer author <author@gmail.com> 2024-09-25T19:32:44+06:30)
					dateParts := strings.SplitN(part, " ", 3)
					date = dateParts[2]
				}

				// Parse parent commit, if it is exist
				if strings.HasPrefix(part, "parent") && part != "" {
					parentParts := strings.Split(part, " ")
					parentCommit = parentParts[1] // Get the parent hash value
				}

				// Parse message content
				if part == "" && len(part) == 0 {
					// Commit message starts after the first empty line
					messageIndex := strings.Index(currentCommitContent, "\n\n")

					if messageIndex != -1 {
						message = currentCommitContent[messageIndex+2:]
					}
				}
			}

			fmt.Printf("commit %s\n", currentCommit)
			fmt.Printf("Author %s\n", author)
			fmt.Printf("Date %s\n", date)
			fmt.Printf("\t%s\n", message)

			// Move to parent commit object
			currentCommit = parentCommit
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
