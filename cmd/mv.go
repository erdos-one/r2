package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// mvCmd represents the mv command
var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "Moves a local file or R2 object to another location locally or in R2.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get client
		c := client("default")

		// If a bucket name is provided, create the bucket
		if len(args) == 2 {
			sourcePath := args[0]
			destinationPath := args[1]
			if !isR2URI(sourcePath) && isR2URI(destinationPath) {
				// Move local file to R2
				destURI := parseR2URI(destinationPath)
				b := c.bucket(destURI.bucket)
				b.upload(sourcePath, destURI.path)
				os.Remove(sourcePath)
			} else if isR2URI(sourcePath) && !isR2URI(destinationPath) {
				// Move R2 object to local file
				sourceURI := parseR2URI(sourcePath)
				b := c.bucket(sourceURI.bucket)
				b.download(sourceURI.path, destinationPath)
				b.delete(sourceURI.path)
			} else if isR2URI(sourcePath) && isR2URI(destinationPath) {
				// Move R2 object to R2 object
				sourceURI := parseR2URI(sourcePath)
				destURI := parseR2URI(destinationPath)
				b := c.bucket(sourceURI.bucket)
				b.copy(sourceURI.path, destURI)
				b.delete(sourceURI.path)
			}
		} else {
			log.Fatal("Please provide both a source and destination path.")
		}
	},
}

func init() {
	// Add the mv subcommand to the root command
	rootCmd.AddCommand(mvCmd)
}
