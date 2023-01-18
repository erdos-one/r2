package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
)

// Convert file size to kb, mb, gb, etc.
func fileSize(b int64) []string {
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

// List all buckets in an account
func listBuckets(client *s3.Client) {
	// Get buckets
	listBucketsOutput, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	// Print creation date and name of each bucket
	for _, object := range listBucketsOutput.Buckets {
		fmt.Println(object.CreationDate.Format("2006-01-02 15:04:05"), *object.Name)
	}
}

// Get all objects in a bucket
func getObjects(client *s3.Client, bucketName string) []types.Object {
	listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Fatal(err)
	}
	return listObjectsOutput.Contents
}

// Get paths of all objects in a bucket
func getObjectPaths(client *s3.Client, bucketName string) []string {
	var objectPaths []string
	for _, object := range getObjects(client, bucketName) {
		objectPaths = append(objectPaths, *object.Key)
	}
	return objectPaths
}

// Print all objects in a bucket
func listObjects(client *s3.Client, bucketName string) {
	// Get creation date, file size, and name of each object
	var objectData [][]string
	for _, object := range getObjects(client, bucketName) {
		// Get file size
		fs := fileSize(object.Size)

		// Append last modified, file size, and file name to objectData
		objectData = append(objectData, []string{
			object.LastModified.Format("2006-01-02 15:04:05"),
			fs[0],
			fs[1],
			*object.Key,
		})
	}

	// Get length of longest file size string
	var longestFileSizeString int
	var longestFileSizeUnitString int
	for _, object := range objectData {
		if len(object[1]) > longestFileSizeString {
			longestFileSizeString = len(object[1])
		}
		if len(object[2]) > longestFileSizeUnitString {
			longestFileSizeUnitString = len(object[2])
		}
	}

	// Print objects
	for _, object := range objectData {
		fmt.Println(
			object[0],
			strings.Repeat(" ", longestFileSizeString-len(object[1])),
			object[1],
			object[2],
			strings.Repeat(" ", longestFileSizeUnitString-len(object[2])),
			object[3],
		)
	}
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List either all buckets or all objects in a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient("default")

		if len(args) > 0 {
			// If args passed to ls, list objects in each bucket passed
			for _, bucketName := range args {
				listObjects(client, bucketName)
			}
		} else {
			// If no args passed to ls, list all buckets
			listBuckets(client)
		}
	},
}

func init() {
	// Add the ls subcommand to the root command
	rootCmd.AddCommand(lsCmd)
}
