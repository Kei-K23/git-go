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
		// Check -m flag is exist
		if commitMessage == "" {
			fmt.Println("No commit message provided. Use -m to provide a message.")
			os.Exit(0) // Exit the command
		}

		// Main logic to handle commit start here

		// Read index file and get blobs to use as tree for project snapshot
		var sb strings.Builder
		entries, err := utils.ReadIndexFile()
		if err != nil {
			log.Fatalln("Error while reading index file content")
		}

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
			log.Fatalln("Error when creating tree object file in .git-go/objects")
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
	},
}

func init() {
	// Register the -m flag
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
