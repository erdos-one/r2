package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/erdos-one/r2/pkg"

	"github.com/spf13/cobra"
)

// Format configuration string
func configString(c pkg.Config) string {
	configTemplate := "[%s]\naccount_id=%s\naccess_key_id=%s\nsecret_access_key=%s"
	return fmt.Sprintf(configTemplate, c.Profile, c.AccountID, c.AccessKeyID, c.SecretAccessKey)
}

// ~/.r2 configuration file path
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(homeDir, ".r2")
}

var R2ConfigFile = getConfigPath()

// Get profile by name or create new one
func getProfile(profileName string) pkg.Config {
	// Get profiles
	profiles := getConfig(false)

	// If profile exists, return it
	for _, profile := range profiles {
		if profile.Profile == profileName {
			return profile
		}
	}

	// Profile doesn't exist, create new one and save to ~/.r2 config file
	profile := getCredentials(profileName)
	writeConfig(profile)

	return profile
}

// Get configuration credentials interactively
func getCredentials(profile string) pkg.Config {
	var c pkg.Config

	// Get profile
	if profile == "" {
		// Get profile name
		fmt.Print("Profile [default]: ")
		fmt.Scanln(&profile)
		if profile == "" {
			profile = "default"
		}
	}
	c.Profile = profile

	// Get account ID
	fmt.Print("Account ID: ")
	fmt.Scanln(&c.AccountID)

	// Get access key ID
	fmt.Print("Access Key ID: ")
	fmt.Scanln(&c.AccessKeyID)

	// Get secret access key
	fmt.Print("Secret Access Key: ")
	fmt.Scanln(&c.SecretAccessKey)

	return c
}

// Parse configuration file and return profiles
func getConfig(createIfNotPresent bool) map[string]pkg.Config {
	// Create configuration file if it doesn't exist
	if _, err := os.Stat(R2ConfigFile); os.IsNotExist(err) {
		// If not creating configuration file, return empty map
		if !createIfNotPresent {
			return make(map[string]pkg.Config)
		}

		f, err := os.Create(R2ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// Get credentials interactively and write to configuration file
		writeConfig(getCredentials(""))
	}

	// Read configuration file
	c, err := os.ReadFile(R2ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Remove empty lines
	configString := regexp.MustCompile(`^\n$`).ReplaceAllString(string(c), "")

	// Parse configuration file into profiles
	var profiles = make(map[string]pkg.Config)

	profilesRe := regexp.MustCompile(`\[[\w\s\]=]+`)
	for _, p := range profilesRe.FindAllString(configString, -1) {
		// Parse profiles
		var profile pkg.Config

		// Get profile name
		if regexp.MustCompile(`\[\w+\]`).MatchString(p) {
			profile.Profile = regexp.MustCompile(`\[(\w+)\]`).FindAllStringSubmatch(p, -1)[0][1]
		}

		// Get account ID
		accountIDRe := regexp.MustCompile(`account_id\s*=\s*(\w+)`)
		if accountIDRe.MatchString(p) {
			profile.AccountID = accountIDRe.FindAllStringSubmatch(p, -1)[0][1]
		}

		// Get access key ID
		akidRe := regexp.MustCompile(`access_key_id\s*=\s*(\w+)`)
		if akidRe.MatchString(p) {
			profile.AccessKeyID = akidRe.FindAllStringSubmatch(p, -1)[0][1]
		}

		// Get secret access key
		sakRe := regexp.MustCompile(`secret_access_key\s*=\s*(\w+)`)
		if sakRe.MatchString(p) {
			profile.SecretAccessKey = sakRe.FindAllStringSubmatch(p, -1)[0][1]
		}

		profiles[profile.Profile] = profile
	}

	return profiles
}

// List all profile names
func listProfiles() []string {
	// Get profiles
	profiles := getConfig(false)

	// Get profile names and sort alphabetically (default profile is always first)
	var profileNames []string
	for _, p := range profiles {
		if p.Profile != "default" {
			profileNames = append(profileNames, p.Profile)
		}
	}

	sort.Slice(profileNames, func(i, j int) bool {
		return strings.ToLower(profileNames[i]) < strings.ToLower(profileNames[j])
	})

	if _, ok := profiles["default"]; ok {
		profileNames = append([]string{"default"}, profileNames...)
	}

	return profileNames
}

// Write configuration to file
func writeConfig(c pkg.Config) {
	// Read configuration file
	profiles := getConfig(false)

	// If not all credentials are provided, fail
	if c.AccountID == "" || c.AccessKeyID == "" || c.SecretAccessKey == "" {
		log.Fatal("All credentials must be provided")
	}

	// Add profile to configuration
	profiles[c.Profile] = c

	// Format profile strings and sort alphabetically (default profile is always first)
	var configStrings []string
	for _, p := range profiles {
		if p.Profile != "default" {
			configStrings = append(configStrings, configString(p))
		}
	}

	sort.Slice(configStrings, func(i, j int) bool {
		return strings.ToLower(configStrings[i]) < strings.ToLower(configStrings[j])
	})

	if _, ok := profiles["default"]; ok {
		configStrings = append([]string{configString(profiles["default"])}, configStrings...)
	}

	// Write configuration to file
	f, err := os.Create(R2ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(strings.Join(configStrings, "\n\n") + "\n")
	if err != nil {
		log.Fatal(err)
	}
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

To list available profiles, run:
  r2 configure --list

To generate an API Token, follow Cloudflare's guide at:
  https://developers.cloudflare.com/r2/data-access/s3-api/tokens/

Be careful not to share your API Token credentials with anyone.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Handle list flag
		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatal(err)
		}
		if list {
			// List profiles
			fmt.Println(strings.Join(listProfiles(), "\n"))
		} else {
			// Parse configuration
			var c pkg.Config
			var err error

			// Get profile name
			c.Profile, err = cmd.Flags().GetString("profile")
			if err != nil {
				log.Fatal(err)
			}

			// Get account ID
			c.AccountID, err = cmd.Flags().GetString("account-id")
			if err != nil {
				log.Fatal(err)
			}

			// Get access key ID
			c.AccessKeyID, err = cmd.Flags().GetString("access-key-id")
			if err != nil {
				log.Fatal(err)
			}

			// Get secret access key
			c.SecretAccessKey, err = cmd.Flags().GetString("secret-access-key")
			if err != nil {
				log.Fatal(err)
			}

			// Either access key ID or secret access key not passed but not both
			if (c.AccessKeyID == "" && c.SecretAccessKey != "") || (c.AccessKeyID != "" && c.SecretAccessKey == "") {
				log.Fatal(`Error: You must either provide both the access key ID and secret access key or
	neither to configure interactively.

	For more information, run:
		r2 help configure`)
			} else {
				// Check if configuration provided
				if c.AccountID != "" && c.AccessKeyID != "" && c.SecretAccessKey != "" {
					writeConfig(c)
				} else {
					// If no configuration provided, get configuration interactively
					writeConfig(getCredentials(""))
				}
			}
		}
	},
}

func init() {
	// Add the configure subcommand to the root command
	rootCmd.AddCommand(configureCmd)

	// Add flags to the configure subcommand
	configureCmd.Flags().BoolP("list", "l", false, "List all named profiles")
	configureCmd.Flags().String("profile", "", "Configure a named profile")
	configureCmd.Flags().String("account-id", "", "R2 Account ID")
	configureCmd.Flags().String("access-key-id", "", "R2 Access Key ID")
	configureCmd.Flags().String("secret-access-key", "", "R2 Secret Access Key")
}
