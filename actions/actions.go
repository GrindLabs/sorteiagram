package actions

import (
	"github.com/ahmdrz/goinsta/v2"

	"github.com/grindlabs/sorteiagram/utils"
	log "github.com/sirupsen/logrus"
)

// ActionsMap - Mapping all available action functions
var ActionsMap = map[string]interface{}{
	"likePost":                LikePost,
	"followProfile":           FollowProfile,
	"followAllProfilesFrom":   FollowAllProfilesFrom,
	"tagFriends":              TagFriends,
	"freeComment":             FreeComment,
	"unfollowProfile":         UnfollowProfile,
	"unfollowAllProfilesFrom": UnfollowAllProfilesFrom,
}

// LikePost - Like a post
func LikePost(instagram *goinsta.Instagram, params ...interface{}) {
	post, err := utils.GetPost(params[0].(string), params[1].(string), instagram)
	logger := log.WithField("post", post.Code)

	if err != nil {
		logger.WithError(err).Panicln("Invalid post code")
	}

	if !post.HasLiked {
		utils.WaitFor(3)

		if err = post.Like(); err != nil {
			logger.WithError(err).Panicln("Unable to like the post")
		}

		logger.Infoln("Post liked successfuly")
		return
	}

	logger.Infoln("Post already liked")
}

// FollowProfile - Follow a single profile
func FollowProfile(instagram *goinsta.Instagram, params ...interface{}) {
	profile := params[0].(string)
	user := utils.GetUser(profile, instagram)
	logger := log.WithField("profile", profile)

	logger.Infoln("Syncing friendship status...")
	utils.WaitFor(3)

	if err := user.FriendShip(); err != nil {
		logger.WithError(err).Warningln("Unable to sync the friendship status")
	}

	if user.Friendship.Following {
		logger.Infoln("Profile already followed")
		return
	}

	utils.WaitFor(3)

	if err := user.Follow(); err != nil {
		logger.WithError(err).Warningln("Unable to follow the profile")
		return
	}

	utils.WaitFor(3)

	if err := user.Mute(goinsta.MuteAll); err != nil {
		logger.WithError(err).Warningln("Unable to mute the profile")
	}

	logger.Infoln("Profile followed successfully")

}

// FollowAllProfilesFrom - Follow all profiles that a profile follows
func FollowAllProfilesFrom(instagram *goinsta.Instagram, params ...interface{}) {
	profile := params[0].(string)
	following := utils.GetUser(profile, instagram).Following()
	logger := log.WithField("from", profile)

	logger.Infoln("Following all profiles...")
	utils.WaitFor(3)

	for following.Next() {
		for _, user := range following.Users {
			logger.Infof("Following @%s...", user.Username)
			FollowProfile(instagram, user.Username)
		}

		utils.WaitFor(3)
	}
}

// UnfollowProfile - Unfollow a single profile
func UnfollowProfile(instagram *goinsta.Instagram, params ...interface{}) {
	profile := params[0].(string)
	user := utils.GetUser(profile, instagram)
	logger := log.WithField("profile", profile)

	logger.Infoln("Syncing friendship status...")
	utils.WaitFor(3)

	if err := user.FriendShip(); err != nil {
		logger.WithError(err).Warningln("Unable to sync the friendship status")
	}

	if !user.Friendship.Following {
		logger.Infoln("Profile already unfollowed")
		return
	}

	utils.WaitFor(3)

	if err := user.Unfollow(); err != nil {
		logger.WithError(err).Warningln("Unable to unfollow the profile")
		return
	}

	logger.Infoln("Profile unfollowed successfully")
}

// UnfollowAllProfilesFrom - Unfollow all profiles that a profile follows
func UnfollowAllProfilesFrom(instagram *goinsta.Instagram, params ...interface{}) {
	profile := params[0].(string)
	following := utils.GetUser(profile, instagram).Following()
	logger := log.WithField("from", profile)

	logger.Infoln("Unfollowing all profiles...")
	utils.WaitFor(3)

	for following.Next() {
		for _, user := range following.Users {
			logger.Infof("Unfollowing @%s...", user.Username)
			UnfollowProfile(instagram, user.Username)
		}

		utils.WaitFor(3)
	}
}

// TagFriends - Tag a pre-defined amount of friends in the post
func TagFriends() {}

// FreeComment - Comment anything in a post
func FreeComment(instagram *goinsta.Instagram, params ...interface{}) {
	post, err := utils.GetPost(params[0].(string), params[1].(string), instagram)
	logger := log.WithFields(log.Fields{
		"profile": params[0].(string),
		"post":    params[1].(string),
		"message": params[2].(string),
	})

	if err != nil {
		logger.WithError(err).Panicln("Invalid post code")
	}

	logger.Infoln("Trying to comment...")
	utils.WaitFor(3)
	post.Comments.Sync()
	utils.WaitFor(3)

	if err = post.Comments.Add(params[2].(string)); err != nil {
		logger.WithError(err).Warningln("Unable to post a comment")
		return
	}

	logger.Infoln("Commented successfully")
}
