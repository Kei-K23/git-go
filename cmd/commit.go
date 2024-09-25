/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Kei-K23/git-go/internal/utils"
	"github.com/spf13/cobra"
)

var commitMessage string // variable to store commit message

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO :: Must TODO
		// Handle and check all updated files are added to staging area and if some file have changes that are different from staged file the commit should not be perform until all file are added to staging area

		// Check -m flag is exist
		if commitMessage == "" {
			fmt.Println("No commit message provided. Use -m to provide a message.")
			os.Exit(0) // Exit the command
		}

		// Main logic to handle commit start here

		// Read index file and get blobs to use as tree for project snapshot
		entries, err := utils.ReadIndexFile()
		if err != nil || len(entries) == 0 {
			fmt.Println("Nothing to commit. Working directory clean.")
			os.Exit(0)
		}

		// Build tree object content
		var sb strings.Builder
		for _, entry := range entries {
			var entryValue string
			// If entry record only have one line, then don't use new line character at the end
			if len(entries) > 1 {
				entryValue = fmt.Sprintf("%s %s %s\n", entry.Mode, entry.Hash, entry.Path)
			} else {
				entryValue = fmt.Sprintf("%s %s %s", entry.Mode, entry.Hash, entry.Path)
			}
			sb.WriteString(entryValue)
		}

		// Generate hash value for tree object
		treeHash, err := utils.HandFileContent([]byte(sb.String()))
		if err != nil {
			log.Fatalf("Cannot generate hash value for given file content")
		}

		// Create tree object and store in objects folder
		os.Mkdir(fmt.Sprintf(".git-go/objects/%s", treeHash[:2]), 0755)
		blogFilePath := fmt.Sprintf(".git-go/objects/%s/%s", treeHash[:2], treeHash[2:])
		treeObj, err := os.Create(blogFilePath)

		if err != nil {
			log.Fatalln("Error when creating tree object")
		}
		defer treeObj.Close()

		// Create new zlib compress writer for tree obj
		var treeCompressBuf bytes.Buffer
		err = utils.CompressContent(&treeCompressBuf, []byte(sb.String()))
		if err != nil {
			log.Fatalln("Error when compressing file content")
		}

		// Add compressed file content call blob to objects file
		treeObj.Write(treeCompressBuf.Bytes())

		// Create commit obj blob
		var commitObjSb strings.Builder
		commitObjSb.WriteString(fmt.Sprintf("tree %s\n", treeHash))

		// Get and check current commit hash value to add as parent commit
		latestCommit := utils.GetCurrentCommit()
		if latestCommit != "" {
			// Add current commit hash value as parent commit when create new commit
			commitObjSb.WriteString(fmt.Sprintf("parent %s\n", latestCommit))
		}

		// Add author and committer info
		// TODO : Handle hard coded value of author name and email
		author := "author <author@gmail.com>" // Replace with dynamic values if needed
		date := utils.GetCurrentTime()        // Use real date-time formatting

		commitObjSb.WriteString(fmt.Sprintf("author %s %s\n", author, date))
		commitObjSb.WriteString(fmt.Sprintf("committer %s %s\n", author, date))
		commitObjSb.WriteString(fmt.Sprintf("\n%s\n", commitMessage))

		commitObjHashValue, err := utils.HandFileContent([]byte(commitObjSb.String()))
		if err != nil {
			log.Fatalf("Cannot generate hash value for commit object")
		}

		// Create commit object and store in objects folder
		os.Mkdir(fmt.Sprintf(".git-go/objects/%s", commitObjHashValue[:2]), 0755)
		commitObjFilePath := fmt.Sprintf(".git-go/objects/%s/%s", commitObjHashValue[:2], commitObjHashValue[2:])
		commitObj, err := os.Create(commitObjFilePath)

		if err != nil {
			log.Fatalln("Error when creating commit object")
		}
		defer commitObj.Close()

		var commitObjCompressBuf bytes.Buffer
		err = utils.CompressContent(&commitObjCompressBuf, []byte(commitObjSb.String()))
		if err != nil {
			log.Fatalln("Error when compressing commit object")
		}
		// Write compressed content to object file
		commitObj.Write(commitObjCompressBuf.Bytes())

		// Add commit hash value as current branch value
		utils.UpdateCommitHashValue(commitObjHashValue)

		fmt.Printf("[master %s] %s\n", commitObjHashValue, commitMessage)
	},
}

func init() {
	// Register the -m flag
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
