package main

// twitter oauth

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/anaxagoras/newsbot/tweet"
	oauth "github.com/araddon/goauth"
	"github.com/araddon/httpstream"
	"html"
	"log"
	"os"
	//"strconv"
	//"strings"
	"encoding/json"
	//"flag"
)

/*var (
	//maxCt    *int    = flag.Int("maxct", 10, "Max # of messages")
	user     *string = flag.String("user", "stinkface", "twitter username")
	ck       *string = flag.String("ck", "MRW1CeiRxufnmR5u3HVObg", "Consumer Key")
	cs       *string = flag.String("cs", "E5XY9qDiZe41rYf5wdpXi6AMNejLUqduqvMdQjCHo", "Consumer Secret")
	ot       *string = flag.String("ot", "14833387-P95MR7O87qZ6pOsuIYU1m1Vo9moaVldLPTv4xXlU2", "Oauth Token")
	osec     *string = flag.String("os", "NdKi8hN2NymcEKQO4JWW02D8cTCbb0ITjCpUtIKnfyYlS", "OAuthTokenSecret")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
	//search   *string = flag.String("search", "android,golang,zeromq,javascript", "keywords to search for, comma delimtted")
	search   *string = flag.String("search", "", "keywords to search for, comma delimtted")
	users    *string = flag.String("users", "1652541,335455570,51241574,18424289,5402612,742143,362051343,1068831", "list of twitter userids to filter for, comma delimtted")
)*/

var config struct {
	LogLevel,
	User,
	ConsumerKey,
	ConsumerSecret,
	OAuthToken,
	OAuthSecret string
	Users    []int64
	Keywords []string
}

func main() {

	if _, err := toml.DecodeFile("newsbot.conf", &config); err != nil {
		fmt.Println(err)
		return
	}

	//flag.Parse()
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), config.LogLevel)

	// make a go channel for sending from listener to processor
	// we buffer it, to help ensure we aren't backing up twitter or else they cut us off
	stream := make(chan []byte, 1000)
	done := make(chan bool)

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
		httpstream.Log(httpstream.ERROR, err.Error())
	} else {

		go func() {
			// while this could be in a different "thread(s)"
			//ct := 0
			var tweet tweet.Tweet
			for tw := range stream {
				//println(string(tw))
				err := json.Unmarshal(tw, &tweet)
				if err != nil {
					httpstream.Log(httpstream.ERROR, err.Error())
				} else {
					tweet.Text = html.UnescapeString(tweet.Text)
					// Tweet parsed
					if userMap[tweet.User.Id] { // If the user is in the list, we're interested
						if tweet.RetweetedStatus.RetweetCount == 0 { // If retweet_count is 0, this is the original author
							fmt.Printf("%s: %s\n", tweet.User.ScreenName, tweet.Text)
						} else { //One of our users is retweeting
							if !userMap[tweet.RetweetedStatus.User.Id] { //this user is not retweeting one of our other users
								fmt.Printf("%s (RT %s): %s\n", tweet.User.ScreenName, tweet.RetweetedStatus.User.ScreenName, tweet.Text)
								println(tweet.Text)
							}
						}
					} else {
						//println("Bad tweet", tweet.Text)
					}
				}

				// heavy lifting
				//ct++
				//if ct > *maxCt {
				//	done <- true
				//}
			}
		}()
		_ = <-done
	}

}
