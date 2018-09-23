package main

import (
    "fmt"
    "time"

    "github.com/SMemsky/go-flakechain/net/p2p"
)

func main() {
    foo, err := p2p.StartNode(12560)
    if err != nil {
        fmt.Println(err)
    }
    defer foo.Stop()

    time.Sleep(0 * time.Second)
}
