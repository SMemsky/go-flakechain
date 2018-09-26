// This package provides basic network P2P functionality
package p2p

import (
	"crypto/rand"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/SMemsky/go-flakechain/net/levin"
)

const (
	maxInConnections  = 64 // TODO
	maxOutConnections = 8

	handshakeInterval                = 60 * time.Second
	connMakerInterval                = 1 * time.Second
	peerlistStoreInterval            = 30 * time.Minute
	grayPeerlistHousekeepingInterval = 1 * time.Minute
	incomingConnectionsInterval      = 15 * time.Minute

	peersPerHandshake = 250

	anchorConnectionsCount      = 2
	whitelistConnectionsPercent = 70
	expectedWhiteConnections    = maxOutConnections * whitelistConnectionsPercent / 100

	handshakeTimeout      = 5 * time.Second
	connectionTimeout     = 5 * time.Second
	pingConnectionTimeout = 2 * time.Second
	invokeTimeout         = 2 * time.Minute

	idlePeerKickTime    = 10 * time.Minute
	passivePeerKickTime = 1 * time.Minute
)

type peerType uint8

const (
	anchorPeer peerType = iota
	whitePeer
	grayPeer
)

var (
	// TODO: Minimum of 12 DNS-resolvable nodes?
	trustedSeedNodes = [...]string{
		"188.35.187.49:12560",
		"188.35.187.51:12560",
		"54.244.21.125:12560",
	}
)

type Node struct {
	// TODO: levin listener
	Ins  []levin.Conn
	Outs []levin.Conn

	whitePeerlist  []PeerListEntry
	grayPeerlist   []PeerListEntry
	anchorPeerlist []AnchorPeerListEntry

	stopIdleRoutine chan struct{}
	wg              sync.WaitGroup
}

// Start runs a node on given port and starts.
// It also runs P2P maintenance routines which should be stopped with Stop
func StartNode(port uint16) (*Node, error) {
	n := &Node{
		Ins:             make([]levin.Conn, 0, maxInConnections),
		Outs:            make([]levin.Conn, 0, maxOutConnections),
		stopIdleRoutine: make(chan struct{}),
	}

	n.wg.Add(1)
	go n.idleRoutine()

	return n, nil
}

// Stop() will block until all open nodes are gracefully closed
func (n *Node) Stop() {
	close(n.stopIdleRoutine)
	n.wg.Wait()
}

func (n *Node) idleRoutine() {
	defer n.wg.Done()

	connMakerTicker := time.NewTicker(connMakerInterval)
	defer connMakerTicker.Stop()

	for {
		select {
		case <-connMakerTicker.C:
			n.makeConnections()
		case <-n.stopIdleRoutine:
			return
		}
	}
}

func (n *Node) makeConnections() {
	log.Println("makeConnections")

	oldConnCount := len(n.Outs)

	if len(n.whitePeerlist) == 0 && len(trustedSeedNodes) != 0 {
		n.connectToSeed()
	}

	if len(n.Outs) < maxOutConnections {
		if len(n.Outs) < expectedWhiteConnections {
			n.makeExpectedConnections(anchorPeer, anchorConnectionsCount)
			n.makeExpectedConnections(whitePeer, expectedWhiteConnections)
			n.makeExpectedConnections(grayPeer, maxOutConnections)
		} else {
			n.makeExpectedConnections(grayPeer, maxOutConnections)
			n.makeExpectedConnections(whitePeer, maxOutConnections)
		}
	}

	if len(n.Outs) == oldConnCount && oldConnCount < maxOutConnections {
		n.connectToSeed()
	}
}

func (n *Node) connectToSeed() {
	if len(trustedSeedNodes) == 0 {
		return
	}

	log.Println("Choosing a seed to connect to")

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(trustedSeedNodes))))
	if err != nil {
		log.Println("connectToSeed:", err)
		return
	}
	n.connectAndHandshakeWithPeer(trustedSeedNodes[index.Int64()], true)
}

func (n *Node) makeExpectedConnections(kind peerType, count uint) {
	log.Println("Trying to meet requirement", count, "for", kind)
}

func (n *Node) connectAndHandshakeWithPeer(address string, onlyTakeThePeers bool) {
	log.Println("Attempting to connect to", address)
}
