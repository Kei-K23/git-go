/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add file contents to the index",
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				// TODO:: Check whether should i use log.Fatalln or fmt.println
				log.Fatalf("File '%s' does not exist", file)
			} else {
				// Read file content that user provided to add to staging area
				fileContentBytes, err := os.ReadFile(file)

				if err != nil {
					log.Fatalf("Cannot read the content of file '%s'", file)
				}
				// Get hash value of file content
				hasher := sha1.New()
				hasher.Write(fileContentBytes)
				hex.EncodeToString(hasher.Sum(nil)) // Get hash value from file content
				// Create file and folder to store compressed blob content

				// Add compressed file content call blob to objects folder
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
