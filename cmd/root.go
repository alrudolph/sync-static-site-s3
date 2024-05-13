package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sync-static-site-s3",
	Short: "Upload a directory containing files for a static site to a S3 Bucket",
	Long: `This CLI uploads the files in a directory to an S3 bucket but makes important
considerations for html files. If the file is an html file, the .html extension is removed
except for index.html and error.html, since we no longer have the html extension we also need
to set the file's Content-Type. This is necessary to allow the files to be accessed without
the .html extension, for example, domain.com/file instead of domain.com/file.html.

Example Usage:
	go run . --directory /path/to/static/site --bucket s3-bucket-name
`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()

		if len(args) > 0 {
			fmt.Println("Additional supplied args will be ignored")
		}

		// Required:
		directory, _ := cmd.Flags().GetString("directory")
		bucket, _ := cmd.Flags().GetString("bucket")

		region, _ := cmd.Flags().GetString("region")

		// Credentials:
		profile, _ := cmd.Flags().GetString("profile")
		accessKeyId, _ := cmd.Flags().GetString("access-key-id")
		secretAccesKey, _ := cmd.Flags().GetString("secret-access-key")

		awsConfig, err := GetAWSConfig(accessKeyId, secretAccesKey, profile, region, ctx)

		if err != nil {
			log.Fatal(err)
		}

		client := s3.NewFromConfig(awsConfig)

		err = EmptyBucket(bucket, client, ctx)

		if err != nil {
			fmt.Println("Failed to clear bucket, aborting upload")
			log.Fatal(err)
		}

		// TODO: can i batch this?
		err = filepath.Walk(
			directory,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				return UploadFile(directory, path, bucket, client, ctx)
			},
		)

		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("directory", "d", "", "Path to the static site directory")
	_ = rootCmd.MarkFlagDirname("directory")
	_ = rootCmd.MarkFlagRequired("directory")

	rootCmd.Flags().StringP("bucket", "b", "", "S3 bucket name")
	_ = rootCmd.MarkFlagRequired("bucket")

	rootCmd.Flags().StringP("region", "r", "us-east-1", "S3 bucket region")

	rootCmd.Flags().String("access-key-id", "", "AWS Access Key ID")
	rootCmd.Flags().String("secret-access-key", "", "AWS Secret Access Key")
	rootCmd.Flags().StringP("profile", "p", "", "AWS Profile name")
}
