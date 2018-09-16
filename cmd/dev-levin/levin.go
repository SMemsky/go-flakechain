package main

import (
	"fmt"
	"encoding/hex"

	"github.com/SMemsky/go-flakechain/net/levin"
)

func main() {
	conn, err := levin.Dial("188.35.187.49:12560")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	packet, _ := hex.DecodeString("01110101010102010108096e6f64655f646174610c100a6c6f63616c5f74696d65051f0f985b00000000076d795f706f727406103100000a6e6574776f726b5f69640a40726e6f77666c616b656e6574776f726b07706565725f696405a1a96dd8d586a7ee0c7061796c6f61645f646174610c101563756d756c61746976655f646966666963756c747905f434072ca80000000e63757272656e745f68656967687405106001000000000006746f705f69640a803baf40bb523e23f437c0b986abd208dca407eb5c46608e51087ec15e995af6ab0b746f705f76657273696f6e0801")
	if err := conn.Invoke(1001, packet); err != nil {
		fmt.Println(err)
	}

	for {
		data, head, err := conn.Receive()
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("%+v\n%x\n", head, data)
	}
}
