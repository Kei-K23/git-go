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

var branchName string // variable to store branch name to use in deletion of branch

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches",
	Run: func(cmd *cobra.Command, args []string) {

		// If delete branch flag exist, then perform branch deletion
		if branchName != "" {
			err := os.Remove(".git-go/refs/heads/" + branchName)
			if err != nil {
				log.Fatalf("Error while deletion of '%s' branch\n", branchName)
			}
			fmt.Printf("Branch name '%s' is deleted\n", branchName)
			os.Exit(0)
		}

		// Create new branch is branch name is provided
		if len(args) > 0 {
			// Check that branch name already exist
			dirEntries, err := os.ReadDir(".git-go/refs/heads")
			if err != nil {
				log.Fatalln("Error while reading heads dir")
			}

			for _, entry := range dirEntries {
				if entry.Name() == args[0] {
					fmt.Printf("Branch '%s' is already exist\n", args[0])
					os.Exit(0)
				}
			}

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
	branchCmd.Flags().StringVarP(&branchName, "delete", "d", "", "Delete the branch")
	rootCmd.AddCommand(branchCmd)
}
