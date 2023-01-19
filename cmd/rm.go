package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an object from an R2 bucket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get client
		c := client("default")

		// If a bucket name is provided, create the bucket
		for _, arg := range args {
			if isR2URI(arg) {
				uri := parseR2URI(arg)
				b := c.bucket(uri.bucket)
				b.delete(uri.path)
			} else {
				log.Fatalf("Path %s is not a valid R2 URI", arg)
			}
		}
	},
}

func init() {
	// Add the rm subcommand to the root command
	rootCmd.AddCommand(rmCmd)
}
