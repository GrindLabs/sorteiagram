package actions

import (
	"time"

	"github.com/ahmdrz/goinsta/v2"

	log "github.com/sirupsen/logrus"
)

// LikePost - Like a post
func LikePost(post, profile string, user *goinsta.User) {
	feed := user.Feed()

	for feed.Next(false) {
		for _, item := range feed.Items {
			if item.Code == post && !item.HasLiked {
				if err := item.Like(); err != nil {
					log.WithError(err).Panicln("Unable to like a post")
				}
			}
		}
	}
}

// FollowProfile - Follow a single profile
func FollowProfile(user *goinsta.User) {
	if err := user.FriendShip(); err != nil {
		log.WithError(err).Panicln("Unable to sync the friendship status")
	}

	if !user.Friendship.Following {
		if err := user.Follow(); err != nil {
			log.WithError(err).Panicln("Unable to follow the profile")
		}
	}
}

// FollowAllProfilesFrom - Follow all profiles that a profile follows
func FollowAllProfilesFrom(user *goinsta.User) {
	following := user.Following()

	for following.Next() {
		for _, u := range following.Users {
			log.Infof("Following @%s...\n", u.Username)
			FollowProfile(&u)
			time.Sleep(10 * time.Second)
		}
	}
}

// TagFriends - Tag a pre-defined amount of friends in the post
func TagFriends(post string, friendsAmount int, isBusiness bool, isVerified bool) {}

// FreeComment - Comment anything in a post
func FreeComment() {}
