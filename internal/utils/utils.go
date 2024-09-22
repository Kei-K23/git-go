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
)

type IndexEntry struct {
	Mode string
	Hash string
	Path string
}

func ReadIndexFile() []IndexEntry {
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

func WriteIndexFile(entries []IndexEntry) error {
	indexFile, err := os.OpenFile(".git-go/index", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	// Close the resource
	defer indexFile.Close()

	for _, entry := range entries {
		entryValue := fmt.Sprintf("%s %s %s\n", entry.Mode, entry.Hash, entry.Path)
		// Write back to index file with updated information
		_, err := indexFile.Write([]byte(entryValue))
		if err != nil {
			return err
		}
	}

	var indexFileContent bytes.Buffer
	indexFile.Read(indexFileContent.Bytes())
	var compressBuf bytes.Buffer

	// Compress the content of index file
	err = CompressContent(&compressBuf, indexFileContent.Bytes())

	if err != nil {
		return err
	}

	// Write compressed content back to index file
	indexFile.Write(compressBuf.Bytes())

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
	defer compressWriter.Close()

	return err
}

// Utility function for file content decompression
func DecompressContent(compressedBuf *bytes.Buffer) ([]byte, error) {
	decompressReader, err := zlib.NewReader(compressedBuf)

	if err != nil {
		return nil, err
	}
	defer decompressReader.Close()
	var decompressBuf bytes.Buffer
	_, err = decompressReader.Read(decompressBuf.Bytes())

	if err != nil {
		return nil, err
	}
	defer decompressReader.Close()

	return decompressBuf.Bytes(), nil
}
