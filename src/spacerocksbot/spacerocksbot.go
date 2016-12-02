// Space Rocks Bot is a bot watching
// asteroids coming too close to earth for the incoming days/week.
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/dns-gh/betterlog"
	conf "github.com/dns-gh/flagsconfig"
	apod "github.com/dns-gh/nasa-apod-client/nasaclient"
	neo "github.com/dns-gh/nasa-neo-client/nasaclient"
	"github.com/dns-gh/twbot"
)

const (
	debugFlag = "debug"
)

// Nasa constants
const (
	// flags definitions
	firstOffsetFlag   = "first-offset"
	offsetFlag        = "offset"
	bodyFlag          = "body"
	pollFrequencyFlag = "poll"
	// TODO save only relevant information on asteroids, the file could become too large at some point otherwise
	nasaPathFlag = "nasa-path"
)

// Twitter constants
const (
	projectURL               = "https://github.com/dns-gh/nasa-space-rocks-bot"
	updateFlag               = "update"
	twitterFollowersPathFlag = "twitter-followers-path"
	twitterFriendsPathFlag   = "twitter-friends-path"
	// TODO save only relevant information on tweets, the file could become too large at some point otherwise
	twitterTweetsPathFlag = "twitter-tweets-path"
	maxFavoriteCountWatch = 2
	maxTryRetweet         = 5
)

var (
	maxRandTimeSleepBetweenRequests = 120 // seconds
)

func main() {
	firstOffset := flag.Int(firstOffsetFlag, 0, "[nasa] offset when fetching data for the first time (days)")
	offset := flag.Int(offsetFlag, 3, "[nasa] offset when fetching data (days)")
	body := flag.String(bodyFlag, "Earth", "[nasa] orbiting body to watch for close asteroids")
	poll := flag.Duration(pollFrequencyFlag, 12*time.Hour, "[nasa] polling frequency of data")
	nasaPath := flag.String(nasaPathFlag, "rocks.json", "[nasa] data file path")
	update := flag.Duration(updateFlag, 30*time.Minute, "[twitter] update frequency of the bot for retweets")
	twitterFollowersPath := flag.String(twitterFollowersPathFlag, "followers.json", "[twitter] data file path for followers")
	twitterFriendsPath := flag.String(twitterFriendsPathFlag, "friends.json", "[twitter] data file path for friends")
	twitterTweetsPath := flag.String(twitterTweetsPathFlag, "tweets.json", "[twitter] data file path for tweets")
	debug := flag.Bool(debugFlag, false, "[twitter] debug mode")
	_, err := conf.NewConfig("nasa.config")
	f, err := betterlog.MakeDateLogger(filepath.Join("Debug", "bot.log"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()
	log.Println("[nasa] first-offset:", *firstOffset)
	log.Println("[nasa] offset:", *offset)
	log.Println("[nasa] body:", *body)
	log.Println("[nasa] poll:", *poll)
	log.Println("[nasa] nasa-path:", *nasaPath)
	log.Println("[twitter] update:", *update)
	log.Println("[twitter] twitter-followers-path:", *twitterFollowersPath)
	log.Println("[twitter] twitter-friends-path:", *twitterFriendsPath)
	log.Println("[twitter] twitter-tweets-path:", *twitterTweetsPath)
	log.Println("[twitter] debug:", *debug)
	bot := twbot.MakeTwitterBot(*twitterFollowersPath, *twitterFriendsPath, *twitterTweetsPath, *debug)
	defer bot.Close()
	neoClient := neo.MakeNasaNeoClient(*firstOffset, *offset, *nasaPath, *body, *debug)
	bot.SetLikePolicy(true, maxFavoriteCountWatch)
	bot.SetRetweetPolicy(maxTryRetweet, true)
	bot.TweetSliceOnceAsync(neoClient.FirstFetch)
	bot.TweetSlicePeriodicallyAsync(neoClient.Fetch, *poll)
	bot.TweetPeriodicallyAsync(func() (string, error) {
		return fmt.Sprintf("check out my source code %s ! 🚀", projectURL), nil
	}, 24*time.Hour)
	bot.RetweetPeriodicallyAsync(searchTweetQueries, bannedQueries, *update)
	policy := &twbot.SleepPolicy{
		MaxRand:               300,
		MaybeSleepChance:      1,
		MaybeSleepTotalChance: 10,
		MaybeSleepMin:         2500,
		MaybeSleepMax:         5000,
	}
	bot.AutoUnfollowFriendsAsync(policy)
	policy.MaybeSleepChance = 5
	policy.MaybeSleepMin = 5000
	policy.MaybeSleepMax = 10000
	bot.AutoFollowFollowersAsync("nasa", 1, policy)
	apodClient := apod.MakeNasaApodClient()
	bot.TweetImagePeriodicallyAsync(apodClient.FetchHD, 24*time.Hour)
	bot.Wait()
}
