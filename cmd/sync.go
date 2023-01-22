package cmd

import (
	"log"

	"github.com/erdos-one/r2/pkg"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs directories and R2 prefixes.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get profile client
		profileName, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}
		c := pkg.Client(getProfile(profileName))

		// If a bucket name is provided, create the bucket
		if len(args) == 2 {
			sourcePath := args[0]
			destinationPath := args[1]
			if !pkg.IsR2URI(sourcePath) && pkg.IsR2URI(destinationPath) {
				// Sync local directory to R2 bucket
				b := c.Bucket(pkg.RemoveR2URIPrefix(destinationPath))
				b.SyncLocalToR2(sourcePath)
			} else if pkg.IsR2URI(sourcePath) && !pkg.IsR2URI(destinationPath) {
				// Sync R2 bucket to local directory
				b := c.Bucket(pkg.RemoveR2URIPrefix(sourcePath))
				b.SyncR2ToLocal(destinationPath)
			} else if pkg.IsR2URI(sourcePath) && pkg.IsR2URI(destinationPath) {
				// Sync R2 bucket to R2 bucket
				b := c.Bucket(pkg.RemoveR2URIPrefix(sourcePath))
				destBucket := c.Bucket(pkg.RemoveR2URIPrefix(destinationPath))
				b.SyncR2ToR2(destBucket)
			}
		} else {
			log.Fatal("Please provide both a source and destination path.")
		}
	},
}

func init() {
	// Add the sync command to the root command
	rootCmd.AddCommand(syncCmd)
}
