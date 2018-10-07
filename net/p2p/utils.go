package p2p

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func randUint64() (result uint64) {
	binary.Read(rand.Reader, binary.LittleEndian, &result)
	return
}
