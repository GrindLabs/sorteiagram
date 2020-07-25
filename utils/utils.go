package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
