# nasa-space-rocks-bot

[![Go Report Card](https://goreportcard.com/badge/github.com/dns-gh/nasa-space-rocks-bot)](https://goreportcard.com/report/github.com/dns-gh/nasa-space-rocks-bot)

The Nasa space rocks bot is a Twitter Bot connected to the Nasa API to get info about threatening asteroids coming near earth.

See https://twitter.com/SpaceRocksBot for one instance of it.

It uses the following nasa client to talk to the Nasa Neo API: https://github.com/dns-gh/nasa-neo-client

It has also a specified and user-defined bot behaviors thanks to https://github.com/dns-gh/twbot.

## Motivation

Simply for fun, practice and trying to do something useful at the same time :)

Still a lot to do! Feel free to join my efforts!

## Installation

- It requires Go language of course. You can set it up by downloading it here: https://golang.org/dl/
- Install it here C:/Go.
- Set your GOPATH, GOROOT and PATH environment variables with:

```
export GOROOT=C:/Go
export GOPATH=WORKING_DIR
export PATH=C:/Go/bin:${PATH}
```

or:

```
@working_dir $ source build/go.sh
```

and then set up your API keys/tokens/secrets:

```
export TWITTER_CONSUMER_KEY="your_twitter_consumer_key"
export TWITTER_CONSUMER_SECRET="your_twitter_consumer_secret"
export TWITTER_ACCESS_TOKEN="your_twitter_access_token"
export TWITTER_ACCESS_SECRET="your_twitter_access_secret"
export NASA_API_KEY="yout_nasa_api_key"
```

You can find get them here: https://api.nasa.gov/index.html#apply-for-an-api-key and https://apps.twitter.com/

## Build and usage

```
@working_dir $ go install nasa-space-rocks-bot
@working_dir $ bin/spacerocksbot.exe -help
  -body string
        [nasa] orbiting body to watch for close asteroids (default "Earth")
  -config string
        configuration filename (default "nasa.config")
  -debug
        [twitter] debug mode
  -first-offset int
        [nasa] offset when fetching data for the first time (days)
  -nasa-path string
        [nasa] data file path (default "rocks.json")
  -offset int
        [nasa] offset when fetching data (days) (default 3)
  -poll duration
        [nasa] polling frequency of data (default 12h0m0s)
  -twitter-followers-path string
        [twitter] data file path for followers (default "followers.json")
  -twitter-friends-path string
        [twitter] data file path for friends (default "friends.json")
  -twitter-tweets-path string
        [twitter] data file path for tweets (default "tweets.json")
  -update duration
        [twitter] update frequency of the bot (default 30m0s)
@working_dir $ bin/spacerocksbot.exe -update=30m -debug=false
[2016-11-18 23:47:34] [info] logging to: D:\WORK\nasa-space-rocks-bot\bin\Debug\bot.log
[2016-11-18 23:47:34] [nasa] first-offset: 3
[2016-11-18 23:47:34] [nasa] offset: 3
[2016-11-18 23:47:34] [nasa] body: Earth
[2016-11-18 23:47:34] [nasa] poll: 12h0m0s
[2016-11-18 23:47:34] [nasa] nasa-path: rocks.json
[2016-11-18 23:47:34] [twitter] update: 30m0s
[2016-11-18 23:47:34] [twitter] twitter-followers-path: followers.json
[2016-11-18 23:47:34] [twitter] twitter-friends-path: friends.json
[2016-11-18 23:47:34] [twitter] twitter-tweets-path: tweets.json
[2016-11-18 23:47:34] [twitter] debug: false
[2016-11-18 23:47:34] [twitter] making twitter bot
[2016-11-18 23:47:34] [nasa] making nasa client
[2016-11-18 23:47:34] [twitter] setting like policy -> auto: true, threshold: 2
[2016-11-18 23:47:34] [twitter] setting retweet policy -> maxTry: 5, like: true
[2016-11-18 23:47:34] [twitter] launching auto follow with 'nasa' over 1 page(s)...
```

## License

See the included LICENSE file.