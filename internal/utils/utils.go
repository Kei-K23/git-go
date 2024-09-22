package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
)

type IndexEntry struct {
	Mode string
	Hash string
	Path string
}

func ReadIndexFile() []IndexEntry {
	// indexFile, err := ÷os(".git-go/index", os.O_APPEND|os.O_WRONLY, 0644)
	var indexFileBuf bytes.Buffer
	indexFile, err := os.Open(".git-go/index")
	if err != nil {
		log.Fatalln("Error when opening index file")
	}
	// Close the resource
	defer indexFile.Close()

	indexFile.Read(indexFileBuf.Bytes())

	// Split the line with next line character
	lines := bytes.Split(indexFileBuf.Bytes(), []byte("\n"))
	var entries []IndexEntry

	for _, line := range lines {
		// Empty line, skip it
		if len(line) == 0 {
			continue
		}

		parts := bytes.Fields(line) // Split by whitespace character
		mode := string(parts[0])
		hash := string(parts[1])
		path := string(parts[2])

		indexEntry := IndexEntry{
			Mode: mode,
			Hash: hash,
			Path: path,
		}
		entries = append(entries, indexEntry) // Add index entry
	}

	return entries
}

func HandFileContent(fileContentBytes []byte) (string, error) {
	// Get hash value of file content
	hasher := sha1.New()
	hasher.Write(fileContentBytes)
	// Get hash value from file content
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
