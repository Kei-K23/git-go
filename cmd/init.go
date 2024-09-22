/*
Copyright © 2024 Kei-K23 <arkar.dev.kei@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Git repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(".git-go")

		// Check .git-go is already exist
		if os.IsExist(err) {
			fmt.Println(".git-go repository already initialized.")
			os.Exit(0) // Success exit
		}

		// Create new .git-go folder with necessary sub-folder and file
		os.Mkdir(".git-go", 0755)
		os.MkdirAll(".git-go/objects", 0755)
		os.MkdirAll(".git-go/refs/heads", 0755)
		os.MkdirAll(".git-go/refs/tags", 0755)

		os.Create(".git-go/config")
		os.Create(".git-go/index")
		headFile, err := os.Create(".git-go/HEAD")

		if err != nil {
			log.Fatalln("Error when creating HEAD file in .git-go")
		}

		content := []byte("ref: refs/heads/master\n")
		headFile.Write(content)
		fmt.Println("Initialized empty .git-go repository.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
