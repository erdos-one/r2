package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

// Remove a R2 bucket
func removeBucket(client *s3.Client, bucketName string) {
	_, err := client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	if err != nil {
		log.Fatalf("Couldn't create bucket %s: %v\n", bucketName, err)
	}
}

// mbCmd represents the mb command
var rbCmd = &cobra.Command{
	Use:   "rb",
	Short: "Remove an R2 bucket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get client
		client := getClient("default")

		// If a bucket name is provided, create the bucket
		if len(args) > 0 {
			bucketName := args[0]
			removeBucket(client, bucketName)
		} else {
			fmt.Println("Please provide a bucket name")
		}
	},
}

func init() {
	// Add the rb subcommand to the root command
	rootCmd.AddCommand(rbCmd)
}
