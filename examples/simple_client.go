package main

import (
	"fmt"
	"log"

	"github.com/gorunriki/mqttc"
)

func main1() {
	client := mqttc.NewClient("localhost:1883", "my-test-client")

	fmt.Println("Connecting to broker...")
	if err := client.Connect(); err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer client.Disconnect()

	fmt.Println("Successfully connected to EMQX!")
	fmt.Println("Press Ctrl+C to disconnect")

	select {}
}
