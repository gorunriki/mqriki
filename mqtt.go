package mqttc

import (
	"errors"
	"net"

	"github.com/gorunriki/mqttc/packets"
)

var (
	ErrNotConnected = errors.New("not connected to broker")
)

type Client struct {
	conn      net.Conn
	broker    string
	clientID  string
	connected bool
}

func NewClien(broker, clientID string) *Client {
	return &Client{
		broker:   broker,
		clientID: clientID,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.broker)
	if err != nil {
		return err
	}
	c.conn = conn

	// create and send CONNECT
	connectPacket := &packets.ConnectPacket{
		ProtocolName:    "MQTT",
		ProtocolVersion: 4,
		CleanSession:    true,
		KeepAlive:       60,
		ClientID:        c.clientID,
	}

	data := packets.EncodeConnect(connectPacket)
	_, err = c.conn.Write(data)
	if err != nil {
		c.conn.Close()
		return err
	}

	// read CONNACK
	resp := make([]byte, 4)
	_, err = c.conn.Read(resp)
	if err != nil {
		c.conn.Close()
		return err
	}

	// verify CONNACK status
	if resp[0] != 0x20 || resp[3] != 0 {
		c.conn.Close()
		return errors.New("connection rejected by broker")
	}

	c.connected = true
	return nil
}

func (c *Client) Disconnect() error {
	if !c.connected {
		return ErrNotConnected
	}

	// send DISCONNECT packet
	disconnectPacket := []byte{0xE0, 0x00}
	c.conn.Write(disconnectPacket)

	c.conn.Close()
	c.connected = false
	return nil
}
