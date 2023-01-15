package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type config struct {
	Profile         string
	AccessKeyID     string
	SecretAccessKey string
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure R2 access",
	Long: `Configure R2 access by providing Cloudflare R2 API Token credentials.

Configuration can be done interactively or by passing flags. If you pass flags,
you must provide both the access key ID and secret access key, otherwise the
command will fail.

To configure interactively, run:
  r2 configure

To configure with flags, run:
  r2 configure --access-key-id <access-key-id> \
    --secret-access-key <secret-access-key>

If you have multiple R2 tokens, you can configure a named profile by passing
the --profile flag.

  Interactively:
    r2 configure --profile my-profile

  With flags:
    r2 configure --profile my-profile --access-key-id <access-key-id> \
      --secret-access-key <secret-access-key>

Profiles are stored in ~/.r2 and can be used by passing the --profile flag to
any command.

To generate an API Token, follow Cloudflare's guide at:
  https://developers.cloudflare.com/r2/data-access/s3-api/tokens/

Be careful not to share your API Token credentials with anyone.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Store configuration options
		var c config

		// Get profile name
		profileName, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatal(err)
		}
		c.Profile = profileName

		// Get access key ID
		accessKeyID, err := cmd.Flags().GetString("access-key-id")
		if err != nil {
			log.Fatal(err)
		}
		c.AccessKeyID = accessKeyID

		// Get secret access key
		secretAccessKey, err := cmd.Flags().GetString("secret-access-key")
		if err != nil {
			log.Fatal(err)
		}
		c.SecretAccessKey = secretAccessKey

		// Parse configuration
		if c == (config{}) {
			// No flags passed, configure interactively
			// Get profile name
			fmt.Print("Profile [default]: ")
			fmt.Scanln(&c.Profile)
			if c.Profile == "" {
				c.Profile = "default"
			}

			// Get access key ID
			fmt.Print("Access Key ID: ")
			fmt.Scanln(&c.AccessKeyID)

			// Get secret access key
			fmt.Print("Secret Access Key: ")
			fmt.Scanln(&c.SecretAccessKey)
		} else if c.Profile != "" && c.AccessKeyID == "" && c.SecretAccessKey == "" {
			// Profile passed, but API credentials not passed
			// Get access key ID
			fmt.Print("Access Key ID: ")
			fmt.Scanln(&c.AccessKeyID)

			// Get secret access key
			fmt.Print("Secret Access Key: ")
			fmt.Scanln(&c.SecretAccessKey)
		} else if c.Profile == "" && c.AccessKeyID != "" && c.SecretAccessKey != "" {
			// Access key ID and secret access key passed, but profile not passed
			// Set profile to default
			c.Profile = "default"
		} else if c.AccessKeyID == "" || c.SecretAccessKey == "" {
			// Either access key ID or secret access key not passed
			log.Fatal(`Error: You must either provide both the access key ID and secret access key or
neither to configure interactively.

For more information, run:
  r2 help configure`)
		}

		fmt.Println(c)
	},
}

func init() {
	// Add the configure subcommand to the root command
	rootCmd.AddCommand(configureCmd)

	// Add flags to the configure subcommand
	configureCmd.Flags().String("profile", "", "Configure a named profile")
	configureCmd.Flags().String("access-key-id", "", "R2 Access Key ID")
	configureCmd.Flags().String("secret-access-key", "", "R2 Secret Access Key")
}
