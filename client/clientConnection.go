package client

import (
	log "github.com/Sirupsen/logrus"
	models "github.com/ruphin/go-carbon/models"
	"regexp"
)

type Client struct {
	subscription models.Subscription
	clientConnection chan string
	// TODO: Stuff that represents remote
}

func (self *Client) run() {
	var dataPoints models.DataPoints
	var ok bool
	var message string
	for {
		select {
		case dataPoints, ok = <- self.subscription.DataChan:
			if !ok {
				log.WithFields(log.Fields{
					"client": self,
				}).Info("Data channel closed")
				self.close()
				break
			}
			log.WithFields(log.Fields{
				"client": self,
				"dataPoints": dataPoints,
			}).Debug("Data Signal")
			self.data(dataPoints)
		case message, ok = <- self.clientConnection:
			if !ok {
				log.WithFields(log.Fields{
					"client": self,
				}).Info("Lost connection to Remote")
				self.subscription.Close()
				break
			}
			log.WithFields(log.Fields{
				"client": self,
				"message": message,
			}).Debug("Message from Remote")
			self.command(message)
		}
	}
}

func NewClient() (client *Client) {
	client = &Client{
		clientConnection: make(chan string),
	}
	go client.run()
	return client
}

var commandRegex = regexp.MustCompile(`^(?P<command>\w+)\s(?P<selector>[\w\.]+)$`)
var commandNames = commandRegex.SubexpNames()

func (self *Client) Test(message string) {
	self.clientConnection <- message
}

func (self *Client) command(message string) {
	r1 := commandRegex.FindAllStringSubmatch(message, -1)
	if r1 == nil {
		log.WithFields(log.Fields{
			"client": self,
			"command": message,
		}).Error("Unknown command from Remote")
		return
	}
	r2 := r1[0]
	md := map[string]string{}
	for i, n := range r2 {
		md[commandNames[i]] = n
	}
	command := md["command"]
	selector := md["selector"]

	switch command {
	case "get":
		self.get(selector)
	case "subscribe":
		self.subscribe(selector)
	}
}

func (self *Client) get(selector string) {
	log.WithFields(log.Fields{
		"selector": selector,
	}).Debug("Get request")
	// TODO: Create subscription and send Get request to Cache
}


func (self *Client) subscribe(selector string) {
	log.WithFields(log.Fields{
		"selector": selector,
	}).Debug("Subscribe request")
	// TODO: Create subscription and send Subscribe request to Cache
}

// Called when the Subscription sends data
func (self *Client) data(dataPoints models.DataPoints) {
	log.WithFields(log.Fields{
		"points": dataPoints,
	}).Debug("New datapoints")
	// TODO: Send dataPoints to remote
}

// Called when the Subscription stops sending data
func (self *Client) close() {
	// TODO: Cleanly close connection to remote
}