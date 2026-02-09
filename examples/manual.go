package main

import (
	"fmt"
	"net"

	"github.com/gorunriki/mqttc/packets"
)

func manual() {
	fmt.Println("Testing manual TCP connection to EMQX...")

	// connect to broker
	conn, err := net.Dial("tcp", "localhost:1883")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create CONNECT packet
	connectPacket := &packets.ConnectPacket{
		ProtocolName:    "MQTT",
		ProtocolVersion: 4,
		CleanSession:    true,
		KeepAlive:       60,
		ClientID:        "test-client-001",
	}

	// encode
	data := packets.EncodeConnect(connectPacket)
	fmt.Printf("Sending %d bytes: %x\n", len(data), data)

	n, err := conn.Write(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sent %d bytes\n", n)

	// read response CONNACK
	response := make([]byte, 1024)
	n, err = conn.Read(response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received %d bytes: %x\n", n, response[:n])

	// parse CONNCACK
	if n >= 4 && response[0] == 0x20 {
		sessionPresent := (response[2] & 0x01) != 0
		returnCode := response[3]

		fmt.Printf("CONNACK received!\n")
		fmt.Printf(" Session Present: %v\n ", sessionPresent)
		fmt.Printf(" Return Code : %d ", returnCode)

		switch returnCode {
		case 0:
			fmt.Println("(Connection Accepted)")
		case 1:
			fmt.Println("(Unacceptable Protocol Version)")
		case 2:
			fmt.Println("(Identifier Rejected)")
		case 3:
			fmt.Println("(Server Unavailable)")
		case 4:
			fmt.Println("(Bad Username/Password)")
		case 5:
			fmt.Println("(Not Authorized)")
		default:
			fmt.Println("(Unknown)")
		}
	} else {
		fmt.Println("Unexpected response")
	}
}
