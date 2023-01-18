// Client-level functions

package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Make new R2 client struct so we can add methods to it
type r2Client struct {
	s3.Client
}

// Get profile by name or create new one
func getProfile(profileName string) config {
	// Get profiles
	profiles := getConfig(false)

	// If profile exists, return it
	for _, profile := range profiles {
		if profile.Profile == profileName {
			return profile
		}
	}

	// Profile doesn't exist, create new one and save to ~/.r2 config file
	profile := getCredentials(profileName)
	writeConfig(profile)

	return profile
}

// Get S3 API Client for given profile
func s3Client(profileName string) *s3.Client {
	// Get profile, if not provided or nonexistent, get configuration interactively
	c := getProfile(profileName)

	// Get R2 account endpoint
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID),
		}, nil
	})

	// Set credentials
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}

// R2 Client for given profile
func client(profileName string) r2Client {
	// Get S3 client to interact with R2
	return r2Client{*s3Client(profileName)}
}

// Make new R2 presign client struct so we can add methods to it
type r2PresignClient struct {
	s3.PresignClient
}

// Get S3 API Client for given profile
func presignClient(profileName string) r2PresignClient {
	s3c := s3Client(profileName)
	return r2PresignClient{*s3.NewPresignClient(s3c)}
}

// List all buckets in an account
func (c *r2Client) listBuckets() {
	// Get buckets
	listBucketsOutput, err := c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	// Print creation date and name of each bucket
	for _, object := range listBucketsOutput.Buckets {
		fmt.Println(object.CreationDate.Format("2006-01-02 15:04:05"), *object.Name)
	}
}

// Create a R2 bucket
func (c *r2Client) createBucket(bucket string) {
	_, err := c.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{},
	})
	if err != nil {
		log.Fatalf("Error creating bucket %s: %v\n", bucket, err)
	}
}

// Remove a R2 bucket
func (c *r2Client) removeBucket(bucket string) {
	_, err := c.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucket)})
	if err != nil {
		log.Fatalf("Couldn't create bucket %s: %v\n", bucket, err)
	}
}
