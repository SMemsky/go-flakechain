// This package provides basic network P2P functionality
package p2p

import (
    "time"

    "github.com/SMemsky/go-flakechain/net/levin"
)

const (
    maxInConnections = 64 // TODO
    maxOutConnections = 8

    handshakeInterval = 60 * time.Second

    peersPerHandshake = 250

    anchorConnectionsCount = 2
    whitelistConnectionsPercent = 70

    handshakeTimeout = 5 * time.Second
    connectionTimeout = 5 * time.Second
    pingConnectionTimeout = 2 * time.Second
    invokeTimeout = 2 * time.Minute
)

var (
    // TODO: Minimum of 12 DNS-resolvable nodes?
    trustedSeedNodes = [...]string {
        "188.35.187.49:12560",
        "188.35.187.51:12560",
        "54.244.21.125:12560",
    }
)

type Node struct {
    // TODO: levin listener
    Ins []levin.Conn
    Outs []levin.Conn
}

// Start runs a node on given port and starts.
// It also runs P2P maintenance routines which should be stopped with Stop
func StartNode(port uint16) (*Node, error) {
    n := &Node{
        Ins: make([]levin.Conn, 0, maxInConnections),
        Outs: make([]levin.Conn, 0, maxOutConnections),
    }

    // go 

    return n, nil
}

func (n *Node) Stop() {
}
