/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"log"
	"os"

	"github.com/Kei-K23/git-go/internal/utils"
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
				// Read file content
				fileContentBytes, err := os.ReadFile(file)
				if err != nil {
					log.Fatalf("Cannot read the content of file '%s'", file)
				}

				// Generate the hash value from file content
				hashValue, err := utils.HandFileContent(fileContentBytes)
				if err != nil {
					log.Fatalf("Cannot read the content of file '%s'", file)
				}

				// Create file and folder to store compressed blob content
				os.Mkdir(fmt.Sprintf(".git-go/objects/%s", hashValue[:2]), 0755)
				blogFilePath := fmt.Sprintf(".git-go/objects/%s/%s", hashValue[:2], hashValue[2:])

				blobFile, err := os.Create(blogFilePath)

				if err != nil {
					log.Fatalln("Error when creating blob file in .git-go/objects")
				}
				defer blobFile.Close()
				// Create new zlib compress writer
				var compressBuf bytes.Buffer
				compressWriter := zlib.NewWriter(&compressBuf)
				_, err = compressWriter.Write(fileContentBytes)

				if err != nil {
					log.Fatalln("Error when compressing file content")
				}

				// Add compressed file content call blob to objects file
				blobFile.Write(compressBuf.Bytes())

				// Add hash value and file name to index file inside .git-go (meaning add file to staging area)
				// example content of index file data without compression 100644 76c7b5bdd3d61bb2657e00b06870f4553294d2a9 README.md
				entries := utils.ReadIndexFile()

				updatedEntires := utils.UpdateIndexFileHashValue(entries, file, hashValue)

				err = utils.WriteIndexFile(updatedEntires)

				if err != nil {
					log.Fatalln("Error when adding entry to index file")
				}

				fmt.Printf("Stored object as : %s\n", blogFilePath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
