# Sync-Static-Site-S3

Empty an S3 Bucket and upload contents of a directory, with file extension and content-type considerations for html files.

## CLI Usage

```
Example Usage:
  go run . --directory /path/to/static/site --bucket s3-bucket-name

Usage:
  sync-static-site-s3 [flags]

Flags:
      --access-key-id string       AWS Access Key ID
  -b, --bucket string              S3 bucket name
  -d, --directory string           Path to the static site directory
  -h, --help                       help for sync-static-site-s3
  -p, --profile string             AWS Profile name
  -r, --region string              S3 bucket region (default "us-east-1")
      --secret-access-key string   AWS Secret Access Key
```

Must use one of the following for credentials:
* Pass in both `access-key-id` and `secret-access-key`
* Pass in `profile`
* Use the environment variable `AWS_PROFILE` as the profile
* Use the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

Download an executable from the [releases](https://github.com/alrudolph/sync-static-site-s3/releases).

## GH Actions Usage

```yaml
- name: Sync Site to S3
  uses: alrudolph/sync-static-site-s3@main
  with:
    bucket: static-site-bucket-name
    directory: path/to/build/folder
  env:
    AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
    AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```
