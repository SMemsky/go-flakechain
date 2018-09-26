package p2p

type BasicNodeData struct {
	LocalTime uint64 `store:"local_time"`
	MyPort    uint32 `store:"my_port"`
	NetworkId string `store:"network_id"`
	PeerId    uint64 `store:"peer_id"`
}

type CoreSyncData struct {
	CumulativeDifficulty uint64 `store:"cumulative_difficulty"`
	CurrentHeight        uint64 `store:"current_height"`
	TopId                string `store:"top_id"`
	TopVersion           uint8  `store:"top_version"`
}

type PeerListEntry struct {
	Address  AddressType `store:"adr"`
	Id       uint64      `store:"id"`
	LastSeen int64       `store:"last_seen"`
}

type AnchorPeerListEntry struct {
	Address   AddressType `store:"adr"`
	Id        uint64      `store:"id"`
	FirstSeen int64       `store:"first_seen"`
}

// Currently Type is always 1. We should add support for IPv6 someday..
type AddressType struct {
	Address struct {
		Ip   uint32 `store:"m_ip"`
		Port uint16 `store:"m_port"`
	} `store:"addr"`
	Type uint8 `store:"type"` // TODO: IPv6 support
}
