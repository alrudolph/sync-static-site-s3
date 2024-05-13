package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func GetAWSConfig(accessKeyId, secretAccessKey, profile, region string, ctx context.Context) (aws.Config, error) {
	// load from profile OR use access key/secret access key (cannot supply both)
	// otherwise, try to use the $AWS_PROFILE profile

	if profile != "" {
		if accessKeyId != "" || secretAccessKey != "" {
			return aws.Config{}, errors.New("cannot provide both profile and access key id/secret access key")
		}
		fmt.Printf("Using profile %s\n", profile)
		c, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))

		if err == nil && c.Region == "" {
			c.Region = region
		}

		return c, err
	}

	if accessKeyId != "" && secretAccessKey != "" {
		fmt.Println("Using access keys")
		return aws.Config{
			Region: region,
			Credentials: credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     accessKeyId,
					SecretAccessKey: secretAccessKey,
				},
			},
		}, nil
	}

	profile = os.Getenv("AWS_PROFILE")

	if profile != "" {
		fmt.Printf("Using default profile %s\n", profile)
		c, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))

		if err == nil && c.Region == "" {
			c.Region = region
		}

		return c, err
	}

	accessKeyId, present := os.LookupEnv("AWS_ACCESS_KEY_ID")

	if !present {
		return aws.Config{}, errors.New("no access key id provided")
	}

	secretAccessKey, present = os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	if !present {
		return aws.Config{}, errors.New("no secret access key provided")
	}

	fmt.Println("Using access keys from environment variables")
	return aws.Config{
		Region: region,
		Credentials: credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKeyId,
				SecretAccessKey: secretAccessKey,
			},
		},
	}, nil
}