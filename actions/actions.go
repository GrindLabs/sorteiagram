package actions

import (
	"errors"
	"time"

	"github.com/ahmdrz/goinsta/v2"

	"github.com/grindlabs/sorteiagram/utils"
	log "github.com/sirupsen/logrus"
)

// ActionsMap - Mapping all available action functions
var ActionsMap = map[string]interface{}{
	"likePost":              LikePost,
	"followProfile":         FollowProfile,
	"followAllProfilesFrom": FollowAllProfilesFrom,
	"tagFriends":            TagFriends,
	"freeComment":           FreeComment,
}

// GetPost - Retrieve a post object
func GetPost(instagram *goinsta.Instagram, params ...interface{}) (goinsta.Item, error) {
	profile := params[0].(string)
	post := params[1].(string)
	feed := utils.GetUser(profile, instagram).Feed()

	log.WithField("post", post).Infoln("Looking for post...")

	for feed.Next(false) {
		for _, item := range feed.Items {
			if item.Code == post {
				log.WithField("post", post).Infoln("Post found")
				return item, nil
			}
		}
	}

	return goinsta.Item{}, errors.New("Post not found")
}

// LikePost - Like a post
func LikePost(instagram *goinsta.Instagram, params ...interface{}) {
	post, err := GetPost(instagram, params...)

	if err != nil {
		log.WithError(err).Panicln("Invalid post code")
	}

	if !post.HasLiked {
		if err = post.Like(); err != nil {
			log.WithError(err).Panicln("Unable to like the post")
		}

		log.WithField("post", post.Code).Infoln("Post liked successfuly")
		return
	}

	log.WithField("post", post.Code).Infoln("Post already liked")
}

// FollowProfile - Follow a single profile
func FollowProfile(instagram *goinsta.Instagram, params ...interface{}) {
	profile := params[0].(string)
	user := utils.GetUser(profile, instagram)

	log.WithField("profile", profile).Infoln("Syncing friendship status...")

	if err := user.FriendShip(); err != nil {
		log.WithError(err).Warningln("Unable to sync the friendship status")
	}

	if !user.Friendship.Following {
		if err := user.Follow(); err != nil {
			log.WithError(err).Warningln("Unable to follow the profile")
			return
		}

		log.WithField("profile", profile).Infoln("Profile followed successfuly")
		return
	}

	log.WithField("profile", profile).Infoln("Profile already followed")
}

// FollowAllProfilesFrom - Follow all profiles that a profile follows
func FollowAllProfilesFrom(instagram *goinsta.Instagram, params ...interface{}) {
	profile := params[0].(string)
	following := utils.GetUser(profile, instagram).Following()

	log.WithField("profile", profile).Infoln("Following all profiles...")

	for following.Next() {
		for _, user := range following.Users {
			log.Infof("Following @%s...\n", user.Username)
			FollowProfile(instagram, user.Username)
			time.Sleep(5 * time.Second)
		}
	}
}

// TagFriends - Tag a pre-defined amount of friends in the post
func TagFriends() {}

// FreeComment - Comment anything in a post
func FreeComment(instagram *goinsta.Instagram, params ...interface{}) {
	post, err := GetPost(instagram, params[0], params[1])

	if err != nil {
		log.WithError(err).Panicln("Invalid post code")
	}

	log.WithField("post", post.Code).Infoln("Trying to comment...")
	post.Comments.Sync()

	if err = post.Comments.Add(params[2].(string)); err != nil {
		log.WithError(err).Warningln("Unable to post a comment")
		return
	}

	log.WithFields(log.Fields{
		"post":    post.Code,
		"message": params[2].(string),
	}).Infoln("Commented successfully")
}
