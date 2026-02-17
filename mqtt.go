package mqttc

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gorunriki/mqttc/packets"
)

var (
	ErrNotConnected = errors.New("not connected to broker")
)

type Client struct {
	conn           net.Conn
	broker         string
	clientID       string
	connected      bool
	messageHandler MessageHandler
	done           chan bool
	incoming       chan *packets.PublishPacket
}

type MessageHandler func(topic string, payload []byte)

func NewClient(broker, clientID string) *Client {
	return &Client{
		broker:   broker,
		clientID: clientID,
		done:     make(chan bool),
		incoming: make(chan *packets.PublishPacket, 100), // buffered channel for incoming messages
	}
}

func (c *Client) SetMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
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

	go c.readLoop()       // start reading incoming packets
	go c.processMessage() // start processing messages
	go c.keepAlive()      // start keep alive pings

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

// needs to change the arguments to add QoS, retail, etc...
func (c *Client) Publish(topic, message string) error {
	if !c.connected {
		return ErrNotConnected
	}

	publishPacket := &packets.PublishPacket{
		Dup:     false,
		QoS:     0,
		Retain:  false,
		Topic:   topic,
		Payload: []byte(message),
	}

	data := packets.EncodePublish(publishPacket)
	_, err := c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Subscribe(topic string) error {
	if !c.connected {
		return ErrNotConnected
	}

	packetID := uint16(1) // or use a counter/generator for multiple subscriptions
	subscriberPacket := &packets.SubscribePacket{
		PacketID: packetID,
		Topics: []packets.Subscription{
			{
				Topic: topic,
				QoS:   0,
			},
		},
	}
	data := packets.EncodeSubscribe(subscriberPacket)
	_, err := c.conn.Write(data)
	if err != nil {
		return err
	}

	// read SUBACK
	resp := make([]byte, 1024) // SUBACK should be at least 5 bytes (fixed header + packet ID + return code)
	n, err := c.conn.Read(resp)
	if err != nil {
		return err
	}

	suback, err := packets.DecodeSuback(resp[:n])
	if err != nil {
		return err
	}

	if suback.PacketID != packetID || len(suback.ReturnCodes) == 0 || suback.ReturnCodes[0] != 0 {
		return errors.New("subscription rejected by broker")
	}

	return nil

}

// function to read incoming packets in a loop
func (c *Client) readLoop() {
	for {
		c.conn.SetReadDeadline(time.Now().Add(45 * time.Second)) // set read timeout to detect disconnections

		resp := make([]byte, 1024)
		_, err := c.conn.Read(resp)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Read timeout, connection may be lost")
				c.connected = false
			} else {
				fmt.Printf("Read error  %v\n", err)
			}
			c.done <- true
			return
		}

		c.conn.SetReadDeadline(time.Time{})

		publish, err := packets.DecodePublish(resp)
		if err != nil {
			fmt.Printf("Publish decode error : %v\n", err)
			continue
		}
		c.incoming <- publish
	}
}

func (c *Client) processMessage() {
	for {
		select {
		case publish := <-c.incoming:
			if c.messageHandler != nil {
				c.messageHandler(publish.Topic, publish.Payload)
			} else {
				fmt.Printf("Received message on topic %s: %s\n", publish.Topic, string(publish.Payload))
			}

			if publish.QoS == 1 {
				c.sendPuback(publish.PacketID)
			}

		case <-c.done:
			return
		}
	}
}

func (c *Client) sendPuback(packetID uint16) error {
	packet := []byte{
		0x40, //PUBACK packet type
		0x02, // remaining length
		byte(packetID >> 8),
		byte(packetID & 0xFF),
	}
	_, err := c.conn.Write(packet)
	return err
}

func (c *Client) keepAlive() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !c.connected {
				return
			}
			fmt.Println("Sending PINGREQ...")
			pingreq := []byte{0xC0, 0x00} // PINGREQ packet
			_, err := c.conn.Write(pingreq)
			if err != nil {
				fmt.Printf("Ping error : %v\n", err)
				c.connected = false
				return
			}
		case <-c.done:
			return
		}
	}
}
