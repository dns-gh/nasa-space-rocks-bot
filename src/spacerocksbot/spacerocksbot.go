// Space Rocks Bot is a bot watching
// asteroids coming too close to earth for the incoming days/week.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	neo "github.com/dns-gh/nasa-neo-client/nasaclient"
	"github.com/dns-gh/twbot"

	conf "github.com/dns-gh/flagsconfig"
)

const (
	debugFlag = "debug"
)

// Nasa constants
const (
	nasaAsteroidsAPIGet = "https://api.nasa.gov/neo/rest/v1/feed?api_key="
	nasaAPIDefaultKey   = "DEMO_KEY"
	nasaTimeFormat      = "2006-01-02"
	fetchMaxSizeError   = "cannot fetch infos for more than 7 days in one request"
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

type timeWriter struct {
	writer io.Writer
}

func (w timeWriter) Write(p []byte) (int, error) {
	date := time.Now().Format("[2006-01-02 15:04:05] ")
	p = append([]byte(date), p...)
	return w.writer.Write(p)
}

func makeDateWriter(w io.Writer) io.Writer {
	return &timeWriter{w}
}

func makeLogger(path string) (string, *os.File, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	err = os.MkdirAll(filepath.Dir(abs), os.ModePerm)
	if err != nil {
		return "", nil, err
	}
	f, err := os.OpenFile(abs, os.O_WRONLY+os.O_APPEND+os.O_CREATE, os.ModePerm)
	return abs, f, err
}

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
	// log to a file also
	log.SetFlags(0)
	logPath, f, err := makeLogger(filepath.Join(filepath.Dir(os.Args[0]), "Debug", "bot.log"))
	if err == nil {
		defer f.Close()
		log.SetOutput(makeDateWriter(io.MultiWriter(f, os.Stderr)))
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("[info] logging to:", logPath)
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
	client := neo.MakeNasaNeoClient(*firstOffset, *offset, *nasaPath, *body, *debug)
	bot.SetLikePolicy(true, maxFavoriteCountWatch)
	bot.SetRetweetPolicy(maxTryRetweet, true)
	bot.TweetSliceOnceAsync(client.FirstFetch)
	bot.TweetSlicePeriodicallyAsync(client.Fetch, *poll)
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
	bot.Wait()
}
