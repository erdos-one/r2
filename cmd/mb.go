package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
)

// Create a R2 bucket
func createBucket(client *s3.Client, bucketName string) {
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{},
	})
	if err != nil {
		log.Fatalf("Error creating bucket %s: %v\n", bucketName, err)
	}
}

// mbCmd represents the mb command
var mbCmd = &cobra.Command{
	Use:   "mb",
	Short: "Create an R2 bucket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get client
		client := getClient("default")

		// If a bucket name is provided, create the bucket
		if len(args) > 0 {
			bucketName := args[0]
			createBucket(client, bucketName)
		} else {
			fmt.Println("Please provide a bucket name")
		}
	},
}

func init() {
	// Add the mb subcommand to the root command
	rootCmd.AddCommand(mbCmd)
}
