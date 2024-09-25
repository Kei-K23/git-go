/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Kei-K23/git-go/internal/utils"
	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches",
	Run: func(cmd *cobra.Command, args []string) {
		// Create new branch is branch name is provided
		if len(args) > 0 {
			path := ".git-go/refs/heads/" + args[0]
			file, err := os.Create(path)
			if err != nil {
				log.Fatalln("Error while creating new branch")
			}
			defer file.Close()

			// Get current commit hash value and added to new created branch
			currentCommit := utils.GetCurrentCommit()

			file.Write([]byte(currentCommit))
			fmt.Printf("New branch call '%s' is created", args[0])
		}

		// No branch name is provided to create, then show all branch name
		dirEntries, err := os.ReadDir(".git-go/refs/heads")
		if err != nil {
			log.Fatalln("Error while reading heads dir")
		}

		for _, entry := range dirEntries {
			fmt.Println(entry.Name())
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
