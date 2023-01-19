package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Check if slice contains string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// Convert file size to kb, mb, gb, etc.
func fileSizeFmt(b int64) []string {
	if b < 1024 {
		// If size is less than 1 KB, return size in bytes
		return []string{fmt.Sprintf("%d", b), "B"}
	} else if b < 1048576 {
		// If size is less than 1 MB, return size in KB
		return []string{fmt.Sprintf("%d", b/1024), "KB"}
	} else if b < 1073741824 {
		// If size is less than 1 GB, return size in MB
		return []string{fmt.Sprintf("%.2f", float64(b)/1048576), "MB"}
	} else {
		// If size is greater than or equal to 1 GB, return size in GB
		return []string{fmt.Sprintf("%.2f", float64(b)/1073741824), "GB"}
	}
}

// Check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Check if a file is a directory
func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.IsDir()
}

// Ensure a directory exists
func ensureDirExists(path string) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
}

// Remove R2 URI prefix
func removeR2URIPrefix(uri string) string {
	return strings.TrimPrefix(uri, "r2://")
}

// Hold R2 URI bucket and file path
type r2URI struct {
	bucket string
	path   string
}

// Determine whether string is an R2 URI
func isR2URI(uri string) bool {
	return strings.HasPrefix(uri, "r2://")
}

// Parse R2 URI
func parseR2URI(uri string) r2URI {
	return r2URI{
		bucket: regexp.MustCompile(`r2://([\w-]+)/.+`).FindStringSubmatch(uri)[1],
		path:   regexp.MustCompile(`r2://[\w-]+/(.+)`).FindStringSubmatch(uri)[1],
	}
}

// Generate MD5 hash of file
func md5sum(path string) string {
	// Get file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Get file hash
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	hashBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashBytes)
}
