// client/examples/chat_demo.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	mqttc "github.com/gorunriki/mqttc"
)

func main() {
	fmt.Println("=== MQTT CHAT DEMO ===")

	// Get user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter chat room: ")
	room, _ := reader.ReadString('\n')
	room = strings.TrimSpace(room)

	// Create client
	clientID := fmt.Sprintf("chat-%s-%d", username, time.Now().Unix())
	client := mqttc.NewClient("localhost:1883", clientID)

	// Connect
	fmt.Println("Connecting to broker...")
	if err := client.Connect(); err != nil {
		panic("Connect failed: " + err.Error())
	}
	defer client.Disconnect()

	fmt.Println("✅ Connected!")

	// Subscribe to chat room
	topic := fmt.Sprintf("chat/%s", room)
	fmt.Printf("Joining chat room: %s\n", topic)

	if err := client.Subscribe(topic); err != nil {
		fmt.Printf("Warning: Could not subscribe: %v\n", err)
	}

	// Goroutine untuk membaca messages
	go func() {
		// TODO: Implement message reading
		fmt.Println("(Message reading not implemented yet)")

	}()

	// Chat loop
	fmt.Println("\nType your messages (Ctrl+C to exit):")
	fmt.Println(strings.Repeat("-", 40))

	for {
		fmt.Printf("[%s] ", username)
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}

		if message == "/exit" {
			break
		}

		// Publish message
		fullMessage := fmt.Sprintf("%s: %s", username, message)
		if err := client.Publish(topic, fullMessage); err != nil {
			fmt.Printf("Error sending: %v\n", err)
		} else {
			fmt.Printf("✓ Sent\n")
		}
	}

	fmt.Println("Goodbye!")
}
