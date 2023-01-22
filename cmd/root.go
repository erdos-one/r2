package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any commands
var rootCmd = &cobra.Command{
	Use:   "r2",
	Short: "Command Line Interface for Cloudflare R2 Storage",
	Long: `r2 is a command line interface for working with Cloudflare's R2 Storage.

Cloudflare's R2 implements the S3 API, attempting to allow users and their
applications to migrate easily, but importantly lacks the key, simple-to-use
features provided by the AWS CLI's s3 subcommand, as opposed to the more complex
and verbose API calls of the s3api subcommand. This CLI fills that gap.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Enable profile flag for all commands
	rootCmd.PersistentFlags().StringP("profile", "p", "default", "R2 profile to use")
}
