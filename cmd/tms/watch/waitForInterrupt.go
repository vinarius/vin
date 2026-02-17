package watch

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func waitForInterrupt(client mqtt.Client) {
	interruptChannel := make(chan os.Signal, 1)

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	<-interruptChannel

	fmt.Println("\nTermination signal received. Shutting down.")

	if client.IsConnected() {
		client.Disconnect(250) // Disconnect with a short grace period
		fmt.Println("MQTT client disconnected.")
	}
}
