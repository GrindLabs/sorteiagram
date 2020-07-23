package cmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Rules - The JSON-object descriptor
type Rules struct {
	LikePostRule              string        `json:"likePostRule"`
	FollowProfileRule         string        `json:"followProfileRule"`
	FollowAllProfilesFromRule string        `json:"followAllProfilesFromRule"`
	TagFriendsRule            []interface{} `json:"tagFriendsRule"`
}

var (
	rules    Rules
	rulesCmd = &cobra.Command{
		Use:   "rules",
		Short: "Create the rules that defines a sweepstakes",
		Long: `Create a set of rules that must be executed to complete the sweepstakes.
Example: sorteiagram rules`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var post string
			var profile string
			var friendsAmount int
			var isBusiness string
			var isVerified string

			fmt.Println("What is the official post URL? ")
			fmt.Scanf("%s", &post)

			fmt.Println("What is the official profile tag? ")
			fmt.Scanf("%s", &profile)

			fmt.Println("How many friends can be tagged by post? ")
			fmt.Scanf("%d", &friendsAmount)

			fmt.Println("Is business profiles allowed to tag? [y/N] ")
			fmt.Scanf("%s", &isBusiness)

			fmt.Println("Is celebrity accounts allowed to tag? [y/N] ")
			fmt.Scanf("%s", &isVerified)

			rules.LikePostRule = post
			rules.FollowProfileRule = profile
			rules.FollowAllProfilesFromRule = profile
			rules.TagFriendsRule = make([]interface{}, 4)
			rules.TagFriendsRule[0] = post
			rules.TagFriendsRule[1] = friendsAmount

			if isBusiness == "y" {
				rules.TagFriendsRule[2] = true
			} else {
				rules.TagFriendsRule[2] = false
			}

			if isVerified == "y" {
				rules.TagFriendsRule[3] = true
			} else {
				rules.TagFriendsRule[3] = false
			}

			hash := sha1.New()
			hash.Write([]byte(post + profile))
			rulesFileName := fmt.Sprintf("%x", hash.Sum(nil))
			jsonData, err := json.Marshal(rules)

			if err != nil {
				log.WithError(err).Panic("Unable to create the rules file")
			}

			path, err := os.Getwd()

			if err != nil {
				log.WithError(err).Panic("Unable to retrieve the rules absolute path")
			}

			ioutil.WriteFile(fmt.Sprintf("%s/rules/%s.json", path, rulesFileName), jsonData, os.ModePerm)
			fmt.Printf("The rules file name is: %s\n", rulesFileName)
		},
	}
)

func init() {
	rootCmd.AddCommand(rulesCmd)
}

// LikePostRule - Like a post
func LikePostRule(post string) {
	fmt.Println(post)
}

// FollowProfileRule - Follow a single profile
func FollowProfileRule(profile string) {}

// FollowAllProfilesFromRule - Follow all profiles that a profile follows
func FollowAllProfilesFromRule(profile string) {}

// TagFriendsRule - Tag a pre-defined amount of friends in the post
func TagFriendsRule(post string, friendsAmount int, isBusiness bool, isVerified bool) {}
