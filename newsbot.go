package main

// twitter oauth

import (
    "log"
)

var config struct {
    Port,
    LogLevel,
    User,
    ConsumerKey,
    ConsumerSecret,
    OAuthToken,
    OAuthSecret string
    Users []int64
    Keywords []string
}

var done =  make(chan bool)

func main() {
    log.Println("Starting newsbot")
    <-done
}
