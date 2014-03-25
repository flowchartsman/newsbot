package main

import (
	"encoding/json"
	"github.com/anaxagoras/newsbot/tweet"
	oauth "github.com/araddon/goauth"
	"github.com/araddon/httpstream"
	"html"
    log "github.com/kdar/factorlog"
    //basiclog "log"
	//"os"
)

func init() {
	// make a go channel for sending from listener to processor
	// we buffer it, to help ensure we aren't backing up twitter or else they cut us off
	stream := make(chan []byte, 1000)

    //TODO: httpstream, your logging is unacceptable.  You will have to go soon
	//httpstream.SetLogger(basiclog.New(os.Stdout, "", basiclog.Ldate|basiclog.Ltime|basiclog.Lshortfile), config.LogLevel)

	httpstream.OauthCon = &oauth.OAuthConsumer{
		Service:          "twitter",
		RequestTokenURL:  "http://twitter.com/oauth/request_token",
		AccessTokenURL:   "http://twitter.com/oauth/access_token",
		AuthorizationURL: "http://twitter.com/oauth/authorize",
		ConsumerKey:      config.ConsumerKey,
		ConsumerSecret:   config.ConsumerSecret,
		CallBackURL:      "oob",
		UserAgent:        "go/newsbotNG",
	}

	//at := goauthcon.GetAccessToken(rt.Token, pin)
	at := oauth.AccessToken{Id: "",
		Token:    config.OAuthToken,
		Secret:   config.OAuthSecret,
		UserRef:  config.User,
		Verifier: "",
		Service:  "twitter",
	}
	// the stream listener effectively operates in one "thread"/goroutine
	// as the httpstream Client processes inside a go routine it opens
	// That includes the handler func we pass in here
	client := httpstream.NewOAuthClient(&at, httpstream.OnlyTweetsFilter(func(line []byte) {
		stream <- line
		// although you can do heavy lifting here, it means you are doing all
		// your work in the same thread as the http streaming/listener
		// by using a go channel, you can send the work to a
		// different thread/goroutine
	}))

	// find list of userids we are going to search for
	//userIds := make([]int64, 0)
	userMap := make(map[int64]bool)
	/*for _, userId := range strings.Split(config.Users, ",") {
			if id, err := strconv.ParseInt(userId, 10, 64); err == nil {
				userIds = append(userIds, id)
	            userMap[id] = true
			}
		}*/

	for _, id := range config.Users {
		userMap[id] = true
	}

	/*var keywords []string
	if search != nil && len(*search) > 0 {
		keywords = strings.Split(config.Keywords, ",")
	}*/
	err := client.Filter(config.Users, config.Keywords, []string{"en"}, nil, false, done)
	if err != nil {
        log.Errorln(err)
	} else {

		go func() {
			var tweet tweet.Tweet

			for {
				select {
				case tw := <-stream:
					//TODO: put unmarshalling, parsing and vetting into another goroutine
					err := json.Unmarshal(tw, &tweet)
					if err != nil {
                        log.Errorln(err)
					} else {
						tweet.Text = html.UnescapeString(tweet.Text)
						// Tweet parsed
						if userMap[tweet.User.Id] { // If the user is in the list, we're interested
							if tweet.RetweetedStatus.RetweetCount == 0 { // If retweet_count is 0, this is the original author
								//println(string(tw))
								log.Debugf("%s: %s %s\n", tweet.User.ScreenName, tweet.Text, tweet.User.ProfileImgURL)
								story := &story{tweet.User.ScreenName, tweet.User.ProfileImgURL, "", tweet.Text}
								messages <- storyMsg(story)
							} else { //One of our users is retweeting
								if !userMap[tweet.RetweetedStatus.User.Id] { //this user is not retweeting one of our other users
									log.Debugf("%s (RT %s): %s\n", tweet.User.ScreenName, tweet.RetweetedStatus.User.ScreenName, tweet.Text)
									println(tweet.Text)
								}
							}
						} else {
							//println("Bad tweet", tweet.Text)
						}
					}
				case <-done:
					break
				}
			}
		}()
	}
}
