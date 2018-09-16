package main

import (
	"fmt"

	"dexter/go-flakechain/net/levin"
	"dexter/go-flakechain/storages/portable"
)

const (
	commandHandshakeId		= 1001
	commandSupportedFlagsId	= 1007
)

type CommandHandshakeRequest struct {
	Node	BasicNodeData	`store:"node_data"`
	Payload	CoreSyncData	`store:"payload_data"`
}

type CommandHandshakeResponse struct {
	// NOTE: Field "local_peerlist" is deprecated and not implemented.
	Peers	[]PeerListEntry	`store:"local_peerlist_new"`
	Node	BasicNodeData	`store:"node_data"`	
}

type BasicNodeData struct {
	LocalTime	uint64	`store:"local_time"`
	MyPort		uint32	`store:"my_port"`
	NetworkId	string	`store:"network_id"`
	PeerId		uint64	`store:"peer_id"`
}

type CoreSyncData struct {
	CumulativeDifficulty	uint64	`store:"cumulative_difficulty"`
	CurrentHeight			uint64	`store:"current_height"`
	TopId					string	`store:"top_id"`
	TopVersion				uint8	`store:"top_version"`
}

type PeerListEntry struct {
	Address		AddressType	`store:"adr"`
	Id			uint64		`store:"id"`
	LastSeen	int64		`store:"last_seen"`
}

type AddressType struct {
	Address	struct {
		Ip		uint32	`store:"m_ip"`
		Port	uint16	`store:"m_port"`
	}					`store:"addr"`
	Type	uint8		`store:"type"` // TODO: IPv6 support
}

func main() {
	conn, err := levin.Dial("188.35.187.49:12560")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	handshake := CommandHandshakeRequest {
		BasicNodeData {
			1536691999,
			12560,
			"rnowflakenetwork",
			17196861955154749857,
		},
		CoreSyncData {
			722293175540,
			90128,
			";\xaf@\xbbR>#\xf47\xc0\xb9\x86\xab\xd2\x08\xdc\xa4\x07\xeb\\F`\x8eQ\x08~\xc1^\x99Z\xf6\xab",
			1,
		},
	}

	rawHandshake, err := portable.Marshal(handshake)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := conn.Invoke(commandHandshakeId, rawHandshake); err != nil {
		fmt.Println(err)
	}

	for {
		data, head, err := conn.Receive()
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("Received packet %d of size %d\n", head.Command, len(data))
		// fmt.Printf("%+v\n%x\n", head, data)
	}
}
