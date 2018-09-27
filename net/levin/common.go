package levin

import (
	"errors"
	"time"
)

const (
	// Appears at the beggining of every levin packet
	levinSignature = 0x0101010101012101

	currentVersion = 1

	bucketSize = 33

	maxPacketSize = 16 * 1024 * 1024 // 16 MiB

	writeTimeout = 10 * time.Second
	readTimeout  = 10 * time.Second

	flagRequest  = 1
	flagResponse = 2
)

var (
	ErrBadSign   = errors.New("net/levin: invalid bucket signature")
	ErrBigPacket = errors.New("net/levin: received packet is too huge")
	ErrVersion   = errors.New("net/levin: packet is of unknown version")

	ErrTimedOut = errors.New("net/levin: operation has timed out")
)

type bucketHead struct {
	Signature       uint64 // Should always be the right one :)
	PacketSize      uint64 // Exactly the size of the data following this header
	ReturnData      bool   // true for INVOKE, false for NOTIFY
	Command         uint32 // Command ID
	ReturnCode      int32  // Always zero?
	Flags           uint32 // 1 - Request, 2 - Response
	ProtocolVersion uint32 // Only version 1 is supported currently
}
