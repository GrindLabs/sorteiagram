package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/ahmdrz/goinsta/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	username        string
	password        string
	sessionFilePath string
	sessionHash     string
	signInCmd       = &cobra.Command{
		Use:   "sign-in",
		Short: "Sign in into an Instagram's account",
		Long: `Use an username and password to sign into an Instagram's account
Example: sorteiagram sign-in USERNAME PASSWORD`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if sessionHash == "" {
				if len(args) > 0 {
					username = args[0]
				}

				if len(args) > 1 {
					password = args[1]
				}

				hash := sha1.New()
				hash.Write([]byte(username + password))
				sessionHash = fmt.Sprintf("%x", hash.Sum(nil))
				log.Debugf("The session hash is %s", sessionHash)
			}

			instagram, err := LoadSession(sessionHash)

			if err != nil {
				log.WithError(err).Warning("Unable to find the session file, trying to signing in...")
				instagram = goinsta.New(username, password)

				if err = instagram.Login(); err != nil {
					log.WithError(err).Panic("Unable to sign into the Instagram's account")
				}
			}

			if err = instagram.Export(sessionFilePath); err != nil {
				log.WithError(err).Panicf("Unable to save the session file in %s", sessionFilePath)
			}

			log.WithFields(log.Fields{
				"username":    instagram.Account.Username,
				"sessionHash": sessionHash,
			}).Infoln("Successfuly sigined in on Instagram")
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&sessionHash, "session-hash", "", "Set a session file (without extension)")
	rootCmd.AddCommand(signInCmd)
}

// LoadSession - Load a session file and return a Instagram object
func LoadSession(sessionHash string) (*goinsta.Instagram, error) {
	path, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	sessionFilePath = fmt.Sprintf("%s/sessions/%s.json", path, sessionHash)
	instagram, err := goinsta.Import(sessionFilePath)
	return instagram, err
}
