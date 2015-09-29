package main

import (
	log "github.com/Sirupsen/logrus"
	client "github.com/ruphin/go-carbon/client"
)

import "time"

func main() {
	var client = client.New()
	log.SetLevel(log.DebugLevel)
	client.Test("get test.something")
	time.Sleep(1 * time.Second)
}
