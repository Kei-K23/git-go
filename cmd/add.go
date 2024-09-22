/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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
				hashValue := hex.EncodeToString(hasher.Sum(nil)) // Get hash value from file content

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

				// TODO:: Make index file handle latest staged hash value for each file
				// Add hash value and file name to index file inside .git-go (meaning add file to staging area)
				// example content of index file data without compression 100644 blob 76c7b5bdd3d61bb2657e00b06870f4553294d2a9    README.md
				indexValue := fmt.Sprintf("100644 %s %s\n", hashValue, file)

				indexFile, err := os.OpenFile(".git-go/index", os.O_APPEND|os.O_WRONLY, 0644)

				if err != nil {
					log.Fatalln("Error when opening index file")
				}

				defer indexFile.Close()

				if _, err := indexFile.Write([]byte(indexValue)); err != nil {
					log.Fatal("Error when writing index file")
				}

				fmt.Printf("Stored object as : %s\n", blogFilePath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
