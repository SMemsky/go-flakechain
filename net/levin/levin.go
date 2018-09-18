package levin

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func Dial(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{conn}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Notify(command uint32, packet []byte) error {
	return c.sendCommand(command, packet, false, 1)
}

func (c *Client) Invoke(command uint32, packet []byte) error {
	return c.sendCommand(command, packet, true, 1)
}

func (c *Client) Respond(command uint32, packet []byte) error {
	return c.sendCommand(command, packet, false, 2)
}

func (c *Client) Receive() ([]byte, *bucketHead, error) {
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

func (c *Client) sendCommand(command uint32, packet []byte, needsReturn bool, flags uint32) error {
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
