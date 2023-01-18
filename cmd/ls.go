package cmd

import (
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List either all buckets or all objects in a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		c := client("default")

		if len(args) > 0 {
			// If args passed to ls, list objects in each bucket passed
			for _, bucketName := range args {
				b := c.bucket(bucketName)
				b.printObjects()
			}
		} else {
			// If no args passed to ls, list all buckets
			c.listBuckets()
		}
	},
}

func init() {
	// Add the ls subcommand to the root command
	rootCmd.AddCommand(lsCmd)
}
