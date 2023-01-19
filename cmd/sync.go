package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs directories and R2 prefixes.",
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
				// Sync local directory to R2 bucket
				b := c.bucket(removeR2URIPrefix(destinationPath))
				b.syncLocalToR2(sourcePath)
			} else if isR2URI(sourcePath) && !isR2URI(destinationPath) {
				// Sync R2 bucket to local directory
				b := c.bucket(removeR2URIPrefix(sourcePath))
				b.syncR2ToLocal(destinationPath)
			} else if isR2URI(sourcePath) && isR2URI(destinationPath) {
				// Sync R2 bucket to R2 bucket
				b := c.bucket(removeR2URIPrefix(sourcePath))
				destBucket := c.bucket(removeR2URIPrefix(destinationPath))
				b.syncR2ToR2(destBucket)
			}
		} else {
			log.Fatal("Please provide both a source and destination path.")
		}
	},
}

func init() {
	// Add the sync subcommand to the root command
	rootCmd.AddCommand(syncCmd)
}
