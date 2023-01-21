// Client-level operations

package pkg

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

type Config struct {
	Profile         string
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
}

// Make new R2 client struct so we can add methods to it
type R2Client struct {
	s3.Client
}

// Get S3 API Client for given profile
func s3Client(c Config) *s3.Client {
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
func Client(c Config) R2Client {
	// Get S3 client to interact with R2
	return R2Client{*s3Client(c)}
}

// Make new R2 presign client struct so we can add methods to it
type R2PresignClient struct {
	s3.PresignClient
}

// Get S3 API Client for given profile
func PresignClient(c Config) R2PresignClient {
	s3c := s3Client(c)
	return R2PresignClient{*s3.NewPresignClient(s3c)}
}

// List all buckets in an account
func (c *R2Client) PrintBuckets() {
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

// Make a R2 bucket
func (c *R2Client) MakeBucket(name string) {
	_, err := c.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{},
	})
	if err != nil {
		log.Fatalf("Error creating bucket %s: %v\n", name, err)
	}
}

// Remove a R2 bucket
func (c *R2Client) RemoveBucket(bucket string) {
	_, err := c.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucket)})
	if err != nil {
		log.Fatalf("Couldn't create bucket %s: %v\n", bucket, err)
	}
}
