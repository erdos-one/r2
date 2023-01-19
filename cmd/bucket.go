// Bucket-level operations

package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

// Put an object in a bucket
func (b *r2Bucket) put(file io.Reader, bucketPath string) error {
	_, err := b.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(bucketPath),
		Body:   file,
	})
	return err
}

// Upload a local file to a bucket
func (b *r2Bucket) upload(localPath, bucketPath string) {
	file, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("Couldn't open file %s to upload: %v\n", localPath, err)
	}

	defer file.Close()

	err = b.put(file, bucketPath)
	if err != nil {
		log.Fatalf("Couldn't upload file %s to r2://%s/%s: %v\n", localPath, b.name, bucketPath, err)
	}
}

// Get an object from a bucket
func (b *r2Bucket) get(bucketPath string) io.ReadCloser {
	obj, err := b.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(bucketPath),
	})
	if err != nil {
		log.Fatalf("Couldn't get file r2://%s/%s: %v\n", b.name, bucketPath, err)
	}

	return obj.Body
}

// Download an object from a bucket
func (b *r2Bucket) download(bucketPath, localPath string) {
	objBody := b.get(bucketPath)

	file, err := os.Create(localPath)
	if err != nil {
		log.Fatalf("Couldn't create file %s to download to: %v\n", localPath, err)
	}

	defer file.Close()
	_, err = io.Copy(file, objBody)
	if err != nil {
		log.Fatalf("Couldn't download file r2://%s/%s to %s: %v\n", b.name, bucketPath, localPath, err)
	}
}

// Copy object from one bucket to another
func (b *r2Bucket) copy(bucketPath string, copyToURI r2URI) {
	_, err := b.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(copyToURI.bucket),
		CopySource: aws.String(b.name + "/" + bucketPath),
		Key:        aws.String(copyToURI.path),
	})
	if err != nil {
		log.Fatalf("Couldn't copy file r2://%s/%s to r2://%s/%s: %v\n", b.name, bucketPath, copyToURI.bucket, copyToURI.path, err)
	}
}

// Delete an object from a bucket
func (b *r2Bucket) delete(bucketPath string) {
	_, err := b.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(bucketPath),
	})
	if err != nil {
		log.Fatalf("Couldn't delete file r2://%s/%s: %v\n", b.name, bucketPath, err)
	}
}

// Sync a local directory to an R2 bucket
func (b *r2Bucket) syncLocalToR2(sourcePath string) {
	// Check if source path exists and is a directory
	if !isDir(sourcePath) {
		log.Fatal("Source path must be a directory.")
	}

	// Get extant paths and their MD5 checksums in bucket
	bucketObjects := make(map[string]string)
	for _, object := range b.getObjects() {
		bucketObjects[*object.Key] = strings.Trim(*object.ETag, `"`)
	}

	// Iterate through paths in source directory
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If path is a file, upload it
		if !info.IsDir() {
			bucketPath := strings.TrimPrefix(path, sourcePath+"/")
			objectMD5, objectInBucket := bucketObjects[bucketPath]
			if !objectInBucket || (md5sum(path) != objectMD5) {
				b.upload(path, bucketPath)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// Sync an R2 bucket to a local directory
func (b *r2Bucket) syncR2ToLocal(destinationPath string) {
	// Check if destination path exists and is a directory
	if !isDir(destinationPath) {
		log.Fatal("Destination path must be a directory.")
	}

	// Iterate through objects and download necessary ones
	for _, object := range b.getObjects() {
		path := *object.Key
		hash := strings.Trim(*object.ETag, `"`)

		// If file either doesn't exist locally or it's changed, download it
		if !fileExists(path) || (fileExists(path) && (md5sum(path) != hash)) {
			outPath := destinationPath + "/" + path
			ensureDirExists(outPath)
			b.download(path, outPath)
		}
	}
}

// Sync an R2 bucket to another R2 bucket
func (b *r2Bucket) syncR2ToR2(destBucket r2Bucket) {
	// Get extant paths and their MD5 checksums in source bucket
	sourceBucketObjects := make(map[string]string)
	for _, object := range b.getObjects() {
		sourceBucketObjects[*object.Key] = strings.Trim(*object.ETag, `"`)
	}

	// Get extant paths and their MD5 checksums in destination bucket
	destBucketObjects := make(map[string]string)
	for _, object := range destBucket.getObjects() {
		destBucketObjects[*object.Key] = strings.Trim(*object.ETag, `"`)
	}

	// Iterate through paths in source bucket and copy necessary ones
	for sourcePath, sourceHash := range sourceBucketObjects {
		destHash, sourceObjectInDestBucket := destBucketObjects[sourcePath]
		if !sourceObjectInDestBucket || (sourceHash != destHash) {
			b.copy(sourcePath, r2URI{bucket: destBucket.name, path: sourcePath})
		}
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
