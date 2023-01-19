package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// presignCmd represents the presign command
var presignCmd = &cobra.Command{
	Use:   "presign",
	Short: "Generate a pre-signed URL for a Cloudflare R2 object",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get profile client
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}
		c := client(profile)
		pc := presignClient(profile)

		for _, arg := range args {
			// Get R2 URI components from argument
			uri := parseR2URI(arg)

			// If object exists in bucket, print presigned URL to get object from bucket, otherwise print
			// presigned URL to put object in bucket
			b := c.bucket(uri.bucket)
			if contains(b.getObjectPaths(), uri.path) {
				fmt.Println(pc.getURL(uri))
			} else {
				fmt.Println(pc.putURL(uri))
			}
		}
	},
}

func init() {
	// Add the presign subcommand to the root command
	rootCmd.AddCommand(presignCmd)
}
