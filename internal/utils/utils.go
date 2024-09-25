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
	"time"
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

// Get current active branch (e.g main or dev, etc...)
func GerCurrentBranch() string {
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
func GetCurrentCommit() string {
	currentBranch := GerCurrentBranch()
	latestCommitFilePath := fmt.Sprintf(".git-go/refs/heads/%s", currentBranch)

	var latestCommitBuf bytes.Buffer
	latestCommitFile, err := os.Open(latestCommitFilePath)
	if os.IsNotExist(err) {
		return ""
	}

	latestCommitBuf.ReadFrom(latestCommitFile)

	// Return the hash value of latest commit obj hash value
	return latestCommitBuf.String()
}

// Function to update the latest commit hash value in the branch reference file
func UpdateCommitHashValue(newHash string) {
	currentBranch := GerCurrentBranch() // Get the current branch (assumed already implemented)
	latestCommitFilePath := fmt.Sprintf(".git-go/refs/heads/%s", currentBranch)

	// Open the file for writing, truncate the content but don't recreate it
	latestCommitFile, err := os.OpenFile(latestCommitFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error when opening current branch file: %v", err)
	}
	defer latestCommitFile.Close()

	// Write the new hash value to the file
	_, err = latestCommitFile.Write([]byte(newHash))
	if err != nil {
		log.Fatalf("Error writing commit hash to file: %v", err)
	}
}

// Read commit object
func ReadCommitObject(hashValue string) []byte {
	var buf bytes.Buffer
	// Path to commit object
	filePath := fmt.Sprintf(".git-go/objects/%s/%s", hashValue[:2], hashValue[2:])

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Error while reading commit object file")
	}

	// Read file content
	buf.ReadFrom(file)
	// Decompress the file content (everything is compress with zlib compression)
	decompressBuf, err := DecompressContent(&buf)
	if err != nil {
		log.Fatalln("Error while decompressing commit object content")
	}

	return decompressBuf // Return decompress bytes buffer
}

// Get hash value of staging file with their file name
func GetHashValueOfStagedFile(filename string) string {
	entries, err := ReadIndexFile()
	if err != nil {
		log.Fatalln("Error while reading index file")
	}

	for _, entry := range entries {
		if entry.Path == filename {
			return entry.Hash
		}
	}
	// Not found, then return empty string
	return ""
}

// Check file is made changes and modified to use in (e.g before adding file to staging area make sure file is modified or not)
func IsFileModified(filename string) bool {
	// Read file content
	fileContentBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Cannot read the content of file '%s'", filename)
	}

	newHashValue, err := HandFileContent(fileContentBytes)
	if err != nil {
		log.Fatalf("Cannot generate hash value for file '%s'", filename)
	}

	oldHashValue := GetHashValueOfStagedFile(filename)

	return newHashValue != oldHashValue
}

// Get the current time in RFC3339 format
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}
