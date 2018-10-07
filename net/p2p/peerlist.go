package p2p

import (
	"sync"
)

const (
	whitePeerlistLimit = 1000
	grayPeerlistLimit  = 5000
)

type peerlist struct {
	mutex sync.Mutex // prevent public access to the lock

	grayPeers   map[string]PeerListEntry
	whitePeers  map[string]PeerListEntry
	anchorPeers map[string]AnchorPeerListEntry
}

func NewPeerlist() *peerlist {
	return &peerlist{
		grayPeers:   make(map[string]PeerListEntry),
		whitePeers:  make(map[string]PeerListEntry),
		anchorPeers: make(map[string]AnchorPeerListEntry),
	}
}

func (p *peerlist) WhiteCount() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.whitePeers)
}

func (p *peerlist) GrayCount() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.grayPeers)
}

func (p *peerlist) MergePeerlist(newPeers []PeerListEntry, localTime int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	peers := make([]PeerListEntry, len(newPeers))
	copy(peers, newPeers)

	if err := fixTimeDelta(peers, localTime); err != nil {
		return
	}

	p.addGrayPeers(peers)
}

func (p *peerlist) GetRandomWhitePeer() (PeerListEntry, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.whitePeers) == 0 {
		return PeerListEntry{}, false
	}

	target := randUint64() % uint64(len(p.whitePeers))
	i := uint64(0)
	for _, v := range p.whitePeers {
		if i == target {
			return v, true
		}
		i++
	}

	return PeerListEntry{}, false
}

func (p *peerlist) GetRandomGrayPeer() (PeerListEntry, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.grayPeers) == 0 {
		return PeerListEntry{}, false
	}

	target := randUint64() % uint64(len(p.grayPeers))
	i := uint64(0)
	for _, v := range p.grayPeers {
		if i == target {
			return v, true
		}
		i++
	}

	return PeerListEntry{}, false
}

func (p *peerlist) addGrayPeers(peers []PeerListEntry) {
	for i := 0; i < len(peers); i++ {
		node := peers[i].Address.String()
		if _, present := p.whitePeers[node]; present {
			continue
		}
		if !isIpAllowed(peers[i].Address.IpString()) {
			continue
		}
		if _, present := p.grayPeers[node]; present {
			// Update only if lastseen is greater
			if p.grayPeers[node].LastSeen < peers[i].LastSeen {
				p.grayPeers[node] = peers[i]
			}
		} else {
			p.grayPeers[node] = peers[i]
		}
	}
}
