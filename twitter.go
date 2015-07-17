package main

import (
	"errors"
	"fmt"
	"github.com/darkhelmet/twitterstream"
	"html"
	"net/url"
	"strconv"
	"strings"
)

type twitterEvent struct {
	tweet *twitterstream.Tweet
	err   error
}

//  TweetStreamer represents a connection
type TweetStreamer struct {
	output     chan *story
	client     *twitterstream.Client
	closeChan  chan struct{}
	connection *twitterstream.Connection
	followList map[string]bool
}

// NewTweetStreamer returns a new twitter streaming API session that's set to
// follow the user IDs specified in the follow slice and output on the provided
// output channel.
func NewTweetStreamer(output chan *story, consumerKey string, consumerSecret string, authToken string, authSecret string, follow []int) (*TweetStreamer, error) {
	//TODO: set UserAgent?
	ts := &TweetStreamer{
		output:     output,
		followList: make(map[string]bool),
		client:     twitterstream.NewClient(config.ConsumerKey, config.ConsumerSecret, config.OAuthToken, config.OAuthSecret),
	}
	for _, id := range follow {
		ts.followList[strconv.Itoa(id)] = true
	}
	err := ts.connect()
	return ts, err
}

// Output returns a read-only channel for retrieving story items
func (t *TweetStreamer) Output() <-chan *story {
	return t.output
}

// Add adds another user to follow (up to 5000). UserIds must be positive.
// Adding the same user ID twice is a no-op.
func (t *TweetStreamer) Add(userID int) error {
	userIDStr := strconv.Itoa(userID)
	switch {
	case userID <= 0:
		return errors.New("Invalid user ID")
	case len(t.followList) == 5000:
		return errors.New("Maximum number of users followed")
	case t.followList[userIDStr]:
		return nil
	}

	t.followList[userIDStr] = true

	//Create the new query string, disconnect, and reconnect with the new query.
	//TODO: Set timer to ensure that reconnections don't fire too often and are
	//bundled
	t.connection.Close()
	close(t.closeChan)
	t.closeChan = make(chan struct{})
	t.connect()
	return nil
}

func (t *TweetStreamer) connect() error {
	if len(t.followList) == 0 {
		return errors.New("No users have been specified or added")
	}

	followIDs := make([]string, 0, len(t.followList))
	for k := range t.followList {
		followIDs = append(followIDs, k)
	}
	var err error
	// url.Values expects map[string][]string to encode multiple values for the
	// same query param, but Twitter expects a comma-separated single value, so:
	t.connection, err = t.client.Filter(url.Values{
		"follow": []string{strings.Join(followIDs, ",")},
		//"track": config.keywords,
	})
	if err != nil {
		return err
	}
	// Set up a buffered channel to receive tweets so that we don't become a slow
	// consumer and get our butts kicked off
	tweetChan := make(chan *twitterEvent, 1000)

	// Goroutine to read tweets and report twitter errors
	// TODO: should we use waitgroups to keep these in line?
	go func() {
		defer close(tweetChan)
	TweetLoop:
		for {
			// Check and see if we need to close
			select {
			case <-t.closeChan:
				break TweetLoop
			default:
			}
			tweet, err := t.connection.Next()
			tweetChan <- &twitterEvent{
				tweet: tweet,
				err:   err,
			}
			if err != nil {
				break TweetLoop
			}
		}
		close(tweetChan)
	}()
	go func() {
		for tw := range tweetChan {
			if tw.err != nil {
				//log.Errorf(err.Error())
			} else {
				tweet := tw.tweet
				fmt.Println("got tweet", tweet)
				// Tweet parsed
				if t.followList[tweet.User.IdStr] { // If the user is in the list, we're interested
					if tweet.RetweetedStatus == nil || //Original user
						tweet.RetweetedStatus.RetweetCount == 0 || //Not been retweeted
						!t.followList[tweet.RetweetedStatus.User.IdStr] { // One of our users retweeting someone else
						tweetText := html.UnescapeString(tweet.Text)
						//println(string(tw))
						//log.Debugf("%s: %s %s\n", tweet.User.ScreenName, tweetText, tweet.User.ProfileImageUrl)
						story := &story{tweet.User.ScreenName, tweet.User.ProfileImageUrl, "", tweetText}
						t.output <- story
					}
				} else {
					//println("Bad tweet", tweet.Text)
				}
			}
		}
	}()
	return nil
}
