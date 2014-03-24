package main

// twitter oauth

import (
	"log"
)

var done = make(chan bool)

func main() {
	log.Println("Starting newsbot")
	<-done
}
