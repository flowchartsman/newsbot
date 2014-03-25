package main

import (
	"github.com/anaxagoras/toml"
	log "github.com/kdar/factorlog"
	"os"
	"path/filepath"
)

var config struct {
	Port,
	LogLevel,
	User,
	ConsumerKey,
	ConsumerSecret,
	OAuthToken,
	OAuthSecret string
	Users    []int64
	Keywords []string
	Scrapers []Scraper
}

var BinPath string

func init() {
	log.SetMinMaxSeverity(log.DEBUG, log.PANIC)
	if _, err := toml.DecodeFile("newsbot.conf", &config); err != nil {
		log.Fatalln(err)
	}

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	BinPath = path

	/*
	   TODO: Sanity checking of config file and directories here. Make sure
	   /static and /template exist, etc.
	*/
}
