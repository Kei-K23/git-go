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
		// Check if the .git-go directory already exists
		_, err := os.Stat(".git-go")
		if !os.IsNotExist(err) {
			// If the .git-go directory exists, exit with a message
			fmt.Println(".git-go repository already initialized.")
			os.Exit(0) // Success exit
		}

		// Create new .git-go folder with necessary sub-folder and files
		err = os.Mkdir(".git-go", 0755)
		if err != nil {
			log.Fatalln("Error creating .git-go directory:", err)
		}

		err = os.MkdirAll(".git-go/objects", 0755)
		if err != nil {
			log.Fatalln("Error creating objects directory:", err)
		}

		err = os.MkdirAll(".git-go/refs/heads", 0755)
		if err != nil {
			log.Fatalln("Error creating heads directory:", err)
		}

		err = os.MkdirAll(".git-go/refs/tags", 0755)
		if err != nil {
			log.Fatalln("Error creating tags directory:", err)
		}

		_, err = os.Create(".git-go/config")
		if err != nil {
			log.Fatalln("Error creating config file:", err)
		}

		_, err = os.Create(".git-go/index")
		if err != nil {
			log.Fatalln("Error creating index file:", err)
		}

		headFile, err := os.Create(".git-go/HEAD")
		if err != nil {
			log.Fatalln("Error when creating HEAD file in .git-go:", err)
		}
		defer headFile.Close()

		// Write the reference to the master branch in the HEAD file
		content := []byte("ref: refs/heads/master\n")
		_, err = headFile.Write(content)
		if err != nil {
			log.Fatalln("Error writing to HEAD file:", err)
		}

		fmt.Println("Initialized empty .git-go repository.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
