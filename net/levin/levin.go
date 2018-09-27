// This package implements basic Levin network functionality
package levin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/SMemsky/go-flakechain/storages/portable"
)

type Conn interface {
	Close()

	// Notify(command uint32, packet []byte) error
	Invoke(commandId uint32, request interface{}, response interface{}, timeout time.Duration) (int32, error)
	// Respond(command uint32, packet []byte) error

	// Receive() ([]byte, *bucketHead, error)

	// Returns a custom, user-defined context
	Context() *interface{}
}

type conn struct {
	conn    net.Conn
	context interface{}

	newMappings chan packetMapping

	stopReceiveRoutine chan struct{}
	wg                 sync.WaitGroup
}

type invokeResponse struct {
	head bucketHead
	data []byte
}

type packetMapping struct {
	id           uint32
	responseChan chan invokeResponse
}

func Dial(address string) (*conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	node := &conn{
		conn:    c,
		context: struct{}{},

		newMappings: make(chan packetMapping),

		stopReceiveRoutine: make(chan struct{}),
	}

	node.wg.Add(1)
	go node.receiveRoutine()
	return node, nil
}

func (c *conn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	close(c.stopReceiveRoutine)
	c.wg.Wait()
}

func (c *conn) Context() *interface{} {
	return &c.context
}

func (c *conn) Invoke(commandId uint32, request interface{}, response interface{}, timeout time.Duration) (int32, error) {
	packet, err := portable.Marshal(request)
	if err != nil {
		return -1, err
	}
	if err = c.sendCommand(commandId, packet, true, flagRequest); err != nil {
		return -1, err
	}

	responseChan := make(chan invokeResponse, 1)
	c.newMappings <- packetMapping{commandId, responseChan}

	select {
	case <-time.After(timeout):
		return -1, ErrTimedOut
	case r, ok := <-responseChan:
		if !ok {
			return -1, fmt.Errorf("Connection closed")
		}
		if r.head.Command != commandId {
			log.Panicln("fixme: please contact devs. This should never happen, lol")
		}

		if err := portable.Unmarshal(r.data, response); err != nil {
			return -1, err
		}
		return r.head.ReturnCode, nil
	}
}

// Receives new packets and directs those where needed. Also listens to
// newMappings channel
func (c *conn) receiveRoutine() {
	defer c.wg.Done()

	head := bucketHead{}
	bucketBuffer := make([]byte, bucketSize)

	responseMap := make(map[uint32](chan invokeResponse))

receiveLoop:
	for {
		select {
		case <-c.stopReceiveRoutine:
			break receiveLoop
		case mapping := <-c.newMappings:
			if _, present := responseMap[mapping.id]; present {
				close(mapping.responseChan)
				break receiveLoop
			}
			responseMap[mapping.id] = mapping.responseChan
		default:
		}

		if _, err := io.ReadFull(c.conn, bucketBuffer); err != nil {
			break receiveLoop
		}
		if err := binary.Read(bytes.NewBuffer(bucketBuffer), binary.LittleEndian, &head); err != nil {
			break receiveLoop
		}

		// // Check response
		if head.Signature != levinSignature {
			break receiveLoop
		}
		if head.ProtocolVersion != currentVersion {
			break receiveLoop
		}
		if head.PacketSize > maxPacketSize {
			break receiveLoop
		}

		data := make([]byte, head.PacketSize)
		if _, err := io.ReadFull(c.conn, data); err != nil {
			break receiveLoop
		}

		if _, present := responseMap[head.Command]; present && head.Flags == flagResponse {
			responseMap[head.Command] <- invokeResponse{head, data}
			close(responseMap[head.Command])
			delete(responseMap, head.Command)
			// log.Println("Sent packet", head.Command, "for handling :)")
		} else {
			// log.Println("Received packet", head.Command, "but did not handle :)")
		}
	}

	// Loop ended, close all invoked sockets
	for _, c := range responseMap {
		close(c)
	}
	// Make SURE there are no new mappings
	done := false
	for !done {
		select {
		case mapping := <-c.newMappings:
			close(mapping.responseChan)
		default:
			done = true
		}
	}
}

func (c *conn) sendCommand(command uint32, packet []byte, needsReturn bool, flags uint32) error {
	head := bucketHead{
		levinSignature,
		uint64(len(packet)),
		needsReturn, command,
		0, flags, 1,
	}

	// c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))

	// Write packet header and data
	if err := binary.Write(c.conn, binary.LittleEndian, head); err != nil {
		return err
	}
	if _, err := c.conn.Write(packet); err != nil {
		return err
	}

	return nil
}
