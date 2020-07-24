package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ahmdrz/goinsta/v2"

	actions "github.com/grindlabs/sorteiagram/actions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the sweepstakes' rules",
	Long: `Start the process of execute the rules that defines a sweepstakes (it needs a logged account).
Example: sorteiagram run RULES_ID`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if sessionHash == "" {
			log.Panicln("Unable to execute the rules because there is no valid session")
		}

		instagram, err := LoadSession(sessionHash)

		if err != nil {
			log.WithError(err).Panicln("Unable to load the Instagram's session")
		}

		path, err := os.Getwd()

		if err != nil {
			log.WithError(err).Panic("Unable to retrieve the rules absolute path")
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/rules/%s.json", path, args[0]))
		var rules map[string]interface{}

		if err = json.Unmarshal(file, &rules); err != nil {
			log.WithError(err).Panic("Unable to load the rules")
		}

		call(rules, instagram)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// Retrieve a User object by their profile name
func getUser(profile string, instagram *goinsta.Instagram) *goinsta.User {
	user, err := instagram.Profiles.ByName(profile)

	if err != nil {
		log.WithError(err).Panicln("Unable to retrieve the user by their profile name")
	}

	return user
}

// Call the respective rules
func call(rules map[string]interface{}, instagram *goinsta.Instagram) {
	for k, v := range rules {
		switch k {
		case "LikePost":
			if v.([]interface{})[0].(string) != "" {
				actions.LikePost(v.([]interface{})[0].(string), v.([]interface{})[1].(string), getUser(v.([]interface{})[1].(string), instagram))
			}
			break

		case "FollowProfile":
			if v.(string) != "" {
				actions.FollowProfile(getUser(v.(string), instagram))
			}
			break

		case "FollowAllProfilesFrom":
			actions.FollowAllProfilesFrom(getUser(v.(string), instagram))
			break

		case "TagFriends":
			break
		}
	}
}
