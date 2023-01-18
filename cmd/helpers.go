package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

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
func getClient(profileName string) *s3.Client {
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

	// Return S3 client to interact with R2
	return s3.NewFromConfig(cfg)
}
