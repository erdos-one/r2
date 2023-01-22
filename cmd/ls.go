package cmd

import (
	"log"

	"github.com/erdos-one/r2/pkg"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List either all buckets or all objects in a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		// Get profile client
		profileName, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}
		c := pkg.Client(getProfile(profileName))

		if len(args) > 0 {
			// If args passed to ls, list objects in each bucket passed
			for _, bucketName := range args {
				// Remove URI scheme if present
				bucketName = pkg.RemoveR2URIPrefix(bucketName)

				b := c.Bucket(bucketName)
				b.PrintObjects()
			}
		} else {
			// If no args passed to ls, list all buckets
			c.PrintBuckets()
		}
	},
}

func init() {
	// Add the ls command to the root command
	rootCmd.AddCommand(lsCmd)
}
