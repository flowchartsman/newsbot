package main

// twitter oauth

import (
    log "github.com/kdar/factorlog"
)

var done = make(chan bool)

func main() {
	log.Infoln("Starting newsbot")
	<-done
}
