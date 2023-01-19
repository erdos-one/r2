package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// mbCmd represents the mb command
var rbCmd = &cobra.Command{
	Use:   "rb",
	Short: "Remove an R2 bucket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get profile client
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}
		c := client(profile)

		// If a bucket name is provided, create the bucket
		if len(args) > 0 {
			c.removeBucket(args[0])
		} else {
			fmt.Println("Please provide a bucket name")
		}
	},
}

func init() {
	// Add the rb subcommand to the root command
	rootCmd.AddCommand(rbCmd)
}
