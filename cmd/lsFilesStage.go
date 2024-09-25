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

// lsFilesStageCmd represents the lsFilesStage command
var lsFilesStageCmd = &cobra.Command{
	Use:   "ls-files-stage",
	Short: "Show information about files in staging area",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := utils.ReadIndexFile()
		if err != nil {
			log.Fatalln("Error while reading index file")
		}

		if len(entries) == 0 {
			fmt.Println("No staging data")
			os.Exit(0)
		}

		for _, entry := range entries {
			fmt.Printf("%s %s %s\n", entry.Mode, entry.Hash, entry.Path)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsFilesStageCmd)
}
