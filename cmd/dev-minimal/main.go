package main

import (
	"log"
	"time"

	"github.com/SMemsky/go-flakechain/net/p2p"
)

const (
	p2pNodeIncomingPort = 12560
)

func main() {
	foo, err := p2p.StartNode(p2pNodeIncomingPort)
	if err != nil {
		log.Println(err)
	}
	defer foo.Stop()

	log.Println("Running!")
	time.Sleep(10 * time.Second)
	log.Println("Ok, gotta stop!")
}
