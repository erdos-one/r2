package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy an object from one R2 path to another",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get profile client
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}
		c := client(profile)

		// If a bucket name is provided, create the bucket
		if len(args) == 2 {
			sourcePath := args[0]
			destinationPath := args[1]
			if !isR2URI(sourcePath) && isR2URI(destinationPath) {
				// Copy local file to R2
				destURI := parseR2URI(destinationPath)
				b := c.bucket(destURI.bucket)
				b.upload(sourcePath, destURI.path)
			} else if isR2URI(sourcePath) && !isR2URI(destinationPath) {
				// Copy R2 object to local file
				sourceURI := parseR2URI(sourcePath)
				b := c.bucket(sourceURI.bucket)
				b.download(sourceURI.path, destinationPath)
			} else if isR2URI(sourcePath) && isR2URI(destinationPath) {
				// Copy R2 object to R2 object
				sourceURI := parseR2URI(sourcePath)
				destURI := parseR2URI(destinationPath)
				b := c.bucket(sourceURI.bucket)
				b.copy(sourceURI.path, destURI)
			}
		} else {
			log.Fatal("Please provide both a source and destination path.")
		}
	},
}

func init() {
	// Add the cp subcommand to the root command
	rootCmd.AddCommand(cpCmd)
}
