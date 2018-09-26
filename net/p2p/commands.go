package p2p

type HandshakeRequest struct {
	NodeData BasicNodeData `store:"node_data"`
	SyncData CoreSyncData  `store:"payload_data"`
}

type HandshakeResponse struct {
	// NOTE: Field "local_peerlist" is deprecated and not implemented.
	// TODO: Remove local_peerlist
	Deprecated string          `store:"local_peerlist"`
	Peers      []PeerListEntry `store:"local_peerlist_new"`
	NodeData   BasicNodeData   `store:"node_data"`
	SyncData   CoreSyncData    `store:"payload_data"`
}

type TimedSyncRequest struct {
	SyncData CoreSyncData `store:"payload_data"`
}

type TimedSyncResponse struct {
	LocalTime  uint64          `store:"local_time"`
	SyncData   CoreSyncData    `store:"payload_data"`
	Deprecated string          `store:"local_peerlist"`
	Peers      []PeerListEntry `store:"local_peerlist_new"`
}

type PingRequest struct {
	// Empty struct
}

type PingResponse struct {
	Status string `store:"status"`
	PeerId uint64 `store:"peer_id"`
}

type SupportedFlagsRequest struct {
	// Empty struct
}

type SupportedFlagsResponse struct {
	Flags uint32 `store:"support_flags"`
}
