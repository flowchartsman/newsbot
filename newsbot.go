package main

// twitter oauth

import (
	log "github.com/kdar/factorlog"
	"os"
)

var done = make(chan bool)

func main() {
	log.SetMinMaxSeverity(log.DEBUG, log.PANIC)
	log.Infoln("Starting newsbot")

	configInit()
	//tweetInit()
	websocketInit()
	webserverInit()

	stories := make(chan *story)

	_, err := NewTweetStreamer(stories, config.ConsumerKey, config.ConsumerSecret, config.OAuthToken, config.OAuthSecret, config.Users)
	if err != nil {
		log.Error("problem initializing tweetstreamer", err)
		os.Exit(1)
	}

	for _, scraper := range config.Scrapers {
		scraper.Output = stories
		scraper.Start()
	}

	go func() {
		for {
			select {
			case webstory := <-stories:
				log.Debugf("Got story %+v", webstory)
				h.broadcast <- storyMsg(webstory)
			}
		}
	}()
	<-done
}
