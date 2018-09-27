// This package implements basic Levin network functionality
package levin

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type Conn interface {
	Close()

	Notify(command uint32, packet []byte) error
	Invoke(command uint32, packet []byte) error
	Respond(command uint32, packet []byte) error

	Receive() ([]byte, *bucketHead, error)

	// Returns a custom, user-defined context
	Context() *interface{}
}

type conn struct {
	conn    net.Conn
	context interface{}
}

func Dial(address string) (*conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &conn{c, struct{}{}}, nil
}

func (c *conn) Close() {
	c.conn.Close()
}

func (c *conn) Context() *interface{} {
	return &c.context
}

func (c *conn) Notify(command uint32, packet []byte) error {
	return c.sendCommand(command, packet, false, 1)
}

func (c *conn) Invoke(command uint32, packet []byte) error {
	return c.sendCommand(command, packet, true, 1)
}

func (c *conn) Respond(command uint32, packet []byte) error {
	return c.sendCommand(command, packet, false, 2)
}

func (c *conn) Receive() ([]byte, *bucketHead, error) {
	head := bucketHead{}

	// TODO: Deadlines
	// c.conn.SetReadDeadline(time.Now().Add(readTimeout))

	// Read response
	buffer := make([]byte, bucketSize)
	if _, err := io.ReadFull(c.conn, buffer); err != nil {
		return nil, nil, err
	}
	if err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &head); err != nil {
		return nil, nil, err
	}

	// // Check response
	if head.Signature != levinSignature {
		return nil, nil, ErrBadSign
	}
	if head.ProtocolVersion != currentVersion {
		return nil, nil, ErrVersion
	}
	if head.PacketSize > maxPacketSize {
		return nil, nil, ErrBigPacket
	}

	response := make([]byte, head.PacketSize)
	if _, err := io.ReadFull(c.conn, response); err != nil {
		return nil, nil, err
	}

	return response, &head, nil
}

func (c *conn) sendCommand(command uint32, packet []byte, needsReturn bool, flags uint32) error {
	head := bucketHead{
		levinSignature,
		uint64(len(packet)),
		needsReturn, command,
		0, flags, 1,
	}

	c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))

	// Write packet header and data
	if err := binary.Write(c.conn, binary.LittleEndian, head); err != nil {
		return err
	}
	if _, err := c.conn.Write(packet); err != nil {
		return err
	}

	return nil
}
