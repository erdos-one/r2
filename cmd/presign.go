package cmd

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

// Get presigned URL for object to get from bucket
func presignGetURL(presignClient *s3.PresignClient, bucketName string, objectName string) string {
	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		log.Fatal("Couldn't get presigned URL for GetObject")
	}
	return presignResult.URL
}

// Get presigned URL for object to put in bucket
func presignPutURL(presignClient *s3.PresignClient, bucketName string, objectName string) string {
	presignResult, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		log.Fatal("Couldn't get presigned URL for PutObject")
	}
	return presignResult.URL
}

// presignCmd represents the presign command
var presignCmd = &cobra.Command{
	Use:   "presign",
	Short: "Generate a pre-signed URL for a Cloudflare R2 object",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get client
		client := getClient("default")

		for _, arg := range args {
			// Parse args
			bucketName := regexp.MustCompile(`r2://([\w-]+)/.+`).FindStringSubmatch(arg)[1]
			objectName := regexp.MustCompile(`r2://[\w-]+/(.+)`).FindStringSubmatch(arg)[1]

			// Create new presign client
			presignClient := s3.NewPresignClient(client)

			// If object exists in bucket, print presigned URL to get object from bucket, otherwise print
			// presigned URL to put object in bucket
			if contains(getObjectPaths(client, bucketName), objectName) {
				fmt.Println(presignGetURL(presignClient, bucketName, objectName))
			} else {
				fmt.Println(presignPutURL(presignClient, bucketName, objectName))
			}
		}
	},
}

func init() {
	// Add the presign subcommand to the root command
	rootCmd.AddCommand(presignCmd)
}
