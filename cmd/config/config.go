package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/alrudolph/snyc-static-site-s3/cmd"
	"github.com/spf13/cobra"
)

func LoadConfigOptions() ([]cmd.SavedConfig, error) {
	usr, err := user.Current()

	if err != nil {
		return nil, err
	}

	homeDir := usr.HomeDir

	configFile, err := os.Open(filepath.Join(homeDir, "sync-s3", "config.json"))

	if err != nil {
		return nil, errors.New("no config profiles found, create one using setup subcommand")
	}

	savedConfig := cmd.SavedConfigFile{}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&savedConfig); err != nil {
		return nil, err
	}

	cwd, err := filepath.Abs(".")

	if err != nil {
		return nil, err
	}

	profiles := []cmd.SavedConfig{}

	for _, profile := range savedConfig.Profiles {
		if profile.UserDirectory == cwd {
			profiles = append(profiles, profile)
			break
		}
	}

	return profiles, nil
}

func starOutWord(word string, showLast int) string {
	if len(word) <= showLast {
		return word
	}

	return fmt.Sprintf("%s%s", word[:showLast], "***")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "See directory configuration",
	// 	Long: `This CLI uploads the files in a directory to an S3 bucket but makes important
	// considerations for html files. If the file is an html file, the .html extension is removed
	// except for index.html and error.html, since we no longer have the html extension we also need
	// to set the file's Content-Type. This is necessary to allow the files to be accessed without
	// the .html extension, for example, domain.com/file instead of domain.com/file.html.

	// Example Usage:
	// 	go run . --directory /path/to/static/site --bucket s3-bucket-name
	// `,
	Run: func(command *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Additional supplied args will be ignored")
		}

		options, err := LoadConfigOptions()

		if err != nil {
			log.Fatal(err)
		}

		if len(options) == 0 {
			log.Fatal("No saved configurations found")
		}

		for _, option := range options {
			fmt.Println(option.Name)

			fmt.Println("    bucket: ", option.Bucket)
			fmt.Println("    region: ", option.Region)
			fmt.Println("    directory: ", option.Directory)

			if option.Profile != "" {
				fmt.Println("    profile: ", option.Profile)
			}

			if option.Role != "" {
				fmt.Println("    role: ", option.Role)
			}

			if option.AccessKeyID != "" {
				fmt.Println("    access key id: ", starOutWord(option.AccessKeyID, 3))
			}

			if option.SecretAccessKey != "" {
				fmt.Println("    secret access key: ", starOutWord(option.SecretAccessKey, 3))
			}

			fmt.Println()
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(configCmd)
}
