// Bucket-level functions

package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type r2Bucket struct {
	client *r2Client
	name   string
}

func (c *r2Client) bucket(bucketName string) r2Bucket {
	return r2Bucket{
		client: c,
		name:   bucketName,
	}
}

// Get all objects in a bucket
func (b *r2Bucket) getObjects() []types.Object {
	listObjectsOutput, err := b.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &b.name,
	})
	if err != nil {
		log.Fatal(err)
	}
	return listObjectsOutput.Contents
}

// Get paths of all objects in a bucket
func (b *r2Bucket) getObjectPaths() []string {
	var objectPaths []string
	for _, object := range b.getObjects() {
		objectPaths = append(objectPaths, *object.Key)
	}
	return objectPaths
}

// Print all objects in a bucket
func (b *r2Bucket) printObjects() {
	// Get creation date, file size, and name of each object
	var objectData [][]string
	for _, object := range b.getObjects() {
		// Get file size
		fs := fileSizeFmt(object.Size)

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

// Get presigned URL for object to get from bucket
func (pc *r2PresignClient) getURL(uri r2URI) string {
	presignResult, err := pc.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(uri.bucket),
		Key:    aws.String(uri.path),
	})
	if err != nil {
		log.Fatal("Couldn't get presigned URL for GetObject")
	}
	return presignResult.URL
}

// Get presigned URL for object to put in bucket
func (pc *r2PresignClient) putURL(uri r2URI) string {
	presignResult, err := pc.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(uri.bucket),
		Key:    aws.String(uri.path),
	})
	if err != nil {
		log.Fatal("Couldn't get presigned URL for PutObject")
	}
	return presignResult.URL
}
