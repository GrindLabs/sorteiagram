package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ahmdrz/goinsta/v2"
	log "github.com/sirupsen/logrus"
)

// GetAbsPath - Retrieve the absolute path
func GetAbsPath() string {
	path, err := os.Getwd()

	if err != nil {
		log.WithError(err).Panic("Unable to retrieve the absolute path")
	}

	return path
}

// LoadSession - Load a session file and return a Instagram object
func LoadSession(sessionHash string) (*goinsta.Instagram, error) {
	sessionFilePath := fmt.Sprintf("%s/sessions/%s.json", GetAbsPath(), sessionHash)
	instagram, err := goinsta.Import(sessionFilePath)

	if err != nil {
		return nil, err
	}

	return instagram, nil
}

// LoadRules - Load the rules data by its id
func LoadRules(rulesID string) map[string][]interface{} {
	jsonFile, err := ioutil.ReadFile(fmt.Sprintf("%s/rules/%s.json", GetAbsPath(), rulesID))

	if err != nil {
		log.WithError(err).Panic("Unable to load the rules file")
	}

	var rules map[string][]interface{}

	if err = json.Unmarshal(jsonFile, &rules); err != nil {
		log.WithError(err).Panic("Unable to parse the rules file")
	}

	return rules
}

// GetUser - Retrieve a User object by their profile name
func GetUser(profile string, instagram *goinsta.Instagram) *goinsta.User {
	user, err := instagram.Profiles.ByName(profile)

	if err != nil {
		log.WithError(err).Panicln("Unable to retrieve the user by their profile name")
	}

	return user
}

// GetPost - Retrieve a post object
func GetPost(profile, post string, instagram *goinsta.Instagram) (goinsta.Item, error) {
	feed := GetUser(profile, instagram).Feed()
	logger := log.WithFields(log.Fields{
		"post":    post,
		"profile": profile,
	})

	logger.Infoln("Looking for post...")
	WaitFor(3)

	for feed.Next(false) {
		for _, item := range feed.Items {
			if item.Code == post {
				logger.Infoln("Post found")
				return item, nil
			}
		}

		WaitFor(3)
	}

	return goinsta.Item{}, errors.New("Post not found")
}

// WaitFor - Wait for a couple of seconds
func WaitFor(seconds time.Duration) {
	time.Sleep(seconds * time.Second)
}
