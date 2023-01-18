package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mbCmd represents the mb command
var mbCmd = &cobra.Command{
	Use:   "mb",
	Short: "Create an R2 bucket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get client
		c := client("default")

		// If a bucket name is provided, create the bucket
		if len(args) > 0 {
			bucketName := args[0]
			c.createBucket(bucketName)
		} else {
			fmt.Println("Please provide a bucket name")
		}
	},
}

func init() {
	// Add the mb subcommand to the root command
	rootCmd.AddCommand(mbCmd)
}
