package p2p

import (
    "net"
)

var (
    disallowedIpRanges []*net.IPNet
)

func isIpAllowed(address string) bool {
    ip := net.ParseIP(address)
    if ip == nil {
        return false
    }

    for _, block := range disallowedIpRanges {
        if block.Contains(ip) {
            return false
        }
    }

    return true
}

func init() {
    for _, blockStr := range []string {
        "127.0.0.0/8",
        "10.0.0.0/8",
        "172.16.0.0/12",
        "192.168.0.0/16",
    } {
        _, block, _ := net.ParseCIDR(blockStr)
        disallowedIpRanges = append(disallowedIpRanges, block)
    }
}
