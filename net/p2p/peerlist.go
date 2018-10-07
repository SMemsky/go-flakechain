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

    grayPeers []PeerListEntry
    whitePeers []PeerListEntry
    anchorPeers []AnchorPeerListEntry
}

func NewPeerlist() *peerlist {
    return peerlist{}
}

func fixTimeDelta(peers []PeerListEntry, localTime int64) error {
    now := time.Now().Unix()
    delta := now - localTime

    for i := 0; i < len(peers); i++ {
        if peers[i].LastSeen > localTime {
            return fmt.Errorf("Peerlist entry from future: local:%d, peer:%d", localTime, peers[i].LastSeen)
        }
        peers[i].LastSeen += delta
    }

    return nil
}

func (p *peerlist) MergePeerlist(newPeers []PeerListEntry, localTime int64) {
    peers := make([]PeerListEntry, len(newPeers))
    copy(peers, newPeers)

    if err := fixTimeDelta(peers, localTime); err != nil {
        return
    }

    p.addGrayPeers(peers)
}

func (p *peerlist) addGrayPeers(peers []PeerListEntry) {
    
}
