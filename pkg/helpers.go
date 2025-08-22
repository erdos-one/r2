package pkg

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

// Contains checks if a string is in a slice of strings.
func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// fileSizeFmt converts a file size in bytes to a human-readable format.
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

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// isDir checks if a path exists and is a directory.
func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	return err == nil && fileInfo.IsDir()
}

// ensureDirExists creates a directory if it does not exist.
func ensureDirExists(path string) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
}

// RemoveR2URIPrefix removes the r2:// prefix from an R2 URI.
func RemoveR2URIPrefix(uri string) string {
	return strings.TrimPrefix(uri, "r2://")
}

// R2URI represents an R2 URI. It contains the bucket name and the path to the file.
type R2URI struct {
	Bucket string
	Path   string
}

// IsR2URI checks if a string is an R2 URI. R2 URI's start with r2://
func IsR2URI(uri string) bool {
	return strings.HasPrefix(uri, "r2://")
}

// ParseR2URISafe parses an R2 URI and returns a R2URI struct.
// Handles URIs with or without paths (e.g., "r2://bucket/" or "r2://bucket/path/").
// Validates that bucket names match the expected format (alphanumeric and hyphens only).
func ParseR2URISafe(uri string) R2URI {
	// Remove the r2:// prefix
	withoutPrefix := strings.TrimPrefix(uri, "r2://")

	// Find the first slash to separate bucket and path
	var bucket, path string
	if idx := strings.Index(withoutPrefix, "/"); idx != -1 {
		bucket = withoutPrefix[:idx]
		path = withoutPrefix[idx+1:] // Everything after the first slash
	} else {
		// No slash found, entire string is bucket name
		bucket = withoutPrefix
		path = ""
	}

	// Validate bucket name format (lowercase alphanumeric and hyphens only)
	// R2 follows S3 naming conventions: lowercase letters, numbers, and hyphens
	if !regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`).MatchString(bucket) {
		log.Fatalf("Invalid bucket name format: %s. Bucket names must contain only lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen.", bucket)
	}

	return R2URI{
		Bucket: bucket,
		Path:   path,
	}
}

// md5sum returns the MD5 hash of a file given its path.
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
