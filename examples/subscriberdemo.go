// client/examples/subscriber_demo.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqttc "github.com/gorunriki/mqttc"
)

func main() {
	fmt.Println("=== MQTT SUBSCRIBER DEMO ===")

	client := mqttc.NewClient("localhost:1883", "subscriber-demo")

	// Set message handler
	client.SetMessageHandler(func(topic string, payload []byte) {
		fmt.Printf("\n📨 [%s] %s\n", topic, string(payload))
	})

	// Connect
	fmt.Println("Connecting to broker...")
	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	fmt.Println("✅ Connected!")

	// Subscribe to multiple topics
	topics := []string{
		"test/#",
		"sensors/temperature/#",
		"chat/#",
	}

	for _, topic := range topics {
		fmt.Printf("Subscribing to: %s\n", topic)
		if err := client.Subscribe(topic); err != nil {
			fmt.Printf("Subscribe error: %v\n", err)
		}
	}

	fmt.Println("\n📡 Waiting for messages... (Ctrl+C to exit)")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Shutting down...")
}
