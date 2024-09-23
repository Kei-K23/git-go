/*
Copyright © 2024 Kei-K23 <arkar.dev.kei@gmail.com>
*/
package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

type IndexEntry struct {
	Mode string
	Hash string
	Path string
}

func ReadIndexFile() ([]IndexEntry, error) {
	// Open the index file
	indexFile, err := os.Open(".git-go/index")
	if err != nil {
		// If the file doesn't exist, return an empty list of entries (new repository)
		if os.IsNotExist(err) {
			return []IndexEntry{}, nil
		}
		log.Fatalln("Error when opening index file:", err)
	}
	defer indexFile.Close()

	// Check if the file is empty
	fileInfo, err := indexFile.Stat()
	if err != nil {
		log.Fatalln("Error when getting file info:", err)
		return nil, err
	}

	// If the file size is zero, return an empty slice
	if fileInfo.Size() == 0 {
		return []IndexEntry{}, nil
	}

	// Read the entire compressed file content into a buffer
	var compressedBuf bytes.Buffer
	_, err = compressedBuf.ReadFrom(indexFile)
	if err != nil {
		return nil, err
	}

	// Decompress the content
	decompressedContent, err := DecompressContent(&compressedBuf)
	if err != nil {
		log.Fatalln("Error when decompressing index file:", err)
		return nil, err
	}

	// Split the decompressed content by lines
	lines := bytes.Split(decompressedContent, []byte("\n"))
	var entries []IndexEntry

	// Parse each line into IndexEntry
	for _, line := range lines {
		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		parts := bytes.Fields(line) // Split by whitespace
		if len(parts) != 3 {
			continue // Ensure the line has 3 parts (mode, hash, path)
		}

		entries = append(entries, IndexEntry{
			Mode: string(parts[0]),
			Hash: string(parts[1]),
			Path: string(parts[2]),
		})
	}

	return entries, nil
}

func WriteIndexFile(entries []IndexEntry) error {
	// Open the file with read/write/truncate mode
	indexFile, err := os.OpenFile(".git-go/index", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// Write entries to the file
	var uncompressedBuf bytes.Buffer
	for _, entry := range entries {
		entryValue := fmt.Sprintf("%s %s %s\n", entry.Mode, entry.Hash, entry.Path)
		uncompressedBuf.Write([]byte(entryValue))
	}

	// Compress the uncompressed content
	var compressBuf bytes.Buffer
	err = CompressContent(&compressBuf, uncompressedBuf.Bytes())
	if err != nil {
		return err
	}

	// Write compressed content to the file
	_, err = indexFile.Write(compressBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func UpdateIndexFileHashValue(entries []IndexEntry, filepath string, newHash string) []IndexEntry {
	var isUpdated bool
	for i, entry := range entries {
		if entry.Path == filepath {
			if entry.Hash != newHash {
				// Update the hash value
				entries[i].Hash = newHash
				isUpdated = true
			}
		}
	}
	// If index entry is new, them add to entry array slice
	if !isUpdated {
		newIndexEntry := IndexEntry{
			Mode: "100644",
			Hash: newHash,
			Path: filepath,
		}
		entries = append(entries, newIndexEntry)
	}

	return entries
}

func HandFileContent(fileContentBytes []byte) (string, error) {
	// Get hash value of file content
	h := sha1.New()
	h.Write(fileContentBytes)
	// Get hash value from file content
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Utility function for file content compression
func CompressContent(compressBuf *bytes.Buffer, fileContentBytes []byte) error {
	compressWriter := zlib.NewWriter(compressBuf)
	_, err := compressWriter.Write(fileContentBytes)
	if err != nil {
		return err
	}
	// Close to ensure that the compressed data is flushed to the buffer
	return compressWriter.Close()
}

// Utility function for file content decompression
func DecompressContent(compressedBuf *bytes.Buffer) ([]byte, error) {
	decompressReader, err := zlib.NewReader(compressedBuf)
	if err != nil {
		return nil, err
	}
	defer decompressReader.Close()

	var decompressBuf bytes.Buffer
	// Use ReadFrom to read decompressed data into the buffer
	_, err = decompressBuf.ReadFrom(decompressReader)
	if err != nil {
		return nil, err
	}

	return decompressBuf.Bytes(), nil
}

// Get current HEAD branch (e.g main or dev, etc...)
func GerCurrentHEAD() string {
	var headFileContentBuf bytes.Buffer
	headFile, err := os.Open(".git-go/HEAD")

	if os.IsNotExist(err) {
		log.Fatalln("HEAD file is not exist")
	}

	// Read head file content
	_, err = headFileContentBuf.ReadFrom(headFile)
	if err != nil {
		log.Fatalln("Error while reading HEAD content")
	}

	parts := strings.Split(headFileContentBuf.String(), " ")
	if len(parts) != 2 || parts[0] != "ref:" {
		log.Fatalln("Invalid HEAD content format")
	}

	// Extract branch name by splitting the ref path by "/"
	refParts := strings.Split(parts[1], "/")
	if len(refParts) < 3 {
		log.Fatalln("Invalid ref format in HEAD file")
	}
	branch := refParts[2] // This is the branch name

	return branch
}

// Get current commit hash
func GetCurrentCommit() {
	currentBranch := GerCurrentHEAD()
	latestCommitFilePath := fmt.Sprintf(".git-go/refs/heads/%s", currentBranch)

	var latestCommitBuf bytes.Buffer
	if _, err := os.Open(latestCommitFilePath); err != nil {
		log.Fatalln("Error while reading commit file")
	}

}
