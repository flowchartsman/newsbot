package main

import (
	"github.com/darkhelmet/twitterstream"
	log "github.com/kdar/factorlog"
	"html"
	"net/url"
	"strconv"
	"strings"
)

type twitterEvent struct {
	tweet *twitterstream.Tweet
	err   error
}

func tweetInit() {
	client := twitterstream.NewClient(config.ConsumerKey, config.ConsumerSecret, config.OAuthToken, config.OAuthSecret)

	//TODO: set UserAgent?

	userMap := make(map[string]bool)

	/*var keywords []string
	if search != nil && len(*search) > 0 {
		keywords = strings.Split(config.Keywords, ",")
	}*/
	toFollow := make([]string, len(config.Users))
	for i, u := range config.Users {
		idStr := strconv.Itoa(u)
		toFollow[i] = idStr
		userMap[idStr] = true
	}
	conn, err := client.Filter(url.Values{
		"follow": []string{strings.Join(toFollow, ",")},
		//track: config.keywords,
	})
	if err != nil {
		log.Errorln(err)
	} else {
		// Set up a buffered channel to receive tweets so that we don't become a slow
		// consumer and get our butts kicked off
		tweetChan := make(chan *twitterEvent, 1000)

		// Goroutine to read tweets and report twitter errors
		go func() {
			for {
				tweet, err := conn.Next()
				tweetChan <- &twitterEvent{
					tweet: tweet,
					err:   err,
				}
				if err != nil {
					close(tweetChan)
					return
				}
			}
		}()
		go func() {
			for t := range tweetChan {
				if t.err != nil {
					log.Errorf(err.Error())
				} else {
					tweet := t.tweet
					log.Debug("got tweet")
					// Tweet parsed
					if userMap[tweet.User.IdStr] { // If the user is in the list, we're interested
						if tweet.RetweetedStatus == nil || //Original user
							tweet.RetweetedStatus.RetweetCount == 0 || //Not been retweeted
							!userMap[tweet.RetweetedStatus.User.IdStr] { // One of our users retweeting someone else
							tweetText := html.UnescapeString(tweet.Text)
							//println(string(tw))
							log.Debugf("%s: %s %s\n", tweet.User.ScreenName, tweetText, tweet.User.ProfileImageUrl)
							story := &story{tweet.User.ScreenName, tweet.User.ProfileImageUrl, "", tweetText}
							h.broadcast <- storyMsg(story)
						}
					} else {
						//println("Bad tweet", tweet.Text)
					}
				}
			}
		}()
	}
}
