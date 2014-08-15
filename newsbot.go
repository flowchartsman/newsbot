package main

// twitter oauth

import (
	log "github.com/kdar/factorlog"
)

var done = make(chan bool)

func main() {
	log.SetMinMaxSeverity(log.DEBUG, log.PANIC)
	log.Infoln("Starting newsbot")

	configInit()
	tweetInit()
	webserverInit()
	websocketInit()

	//TODO: move the twitter handling code in here, too
	stories := make(chan *story)

	for _, scraper := range config.Scrapers {
		scraper.Output = stories
		scraper.Start()
	}

	go func() {
		for {
			select {
			case webstory := <-stories:
				log.Debugf("Got story %+v", webstory)
				messages <- storyMsg(webstory)
			}
		}
	}()
	<-done
}
