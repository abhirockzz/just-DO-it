package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/nats-io/go-nats"
)

const subject = "foo"

func main() {

	natsURL := os.Getenv("NATS_SERVER")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	fmt.Println("Connecting to NATS server - " + natsURL)

	opts := nats.Options{
		AllowReconnect: true,
		MaxReconnect:   5,
		ReconnectWait:  5 * time.Second,
		Timeout:        3 * time.Second,
		Url:            natsURL,
	}

	conn, _ := opts.Connect()
	fmt.Println("Subscriber discovered servers - ", conn.DiscoveredServers())
	fmt.Println("Subscriber connected to NATS server with ID - ", conn.ConnectedServerId())

	conn.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Println("Got message " + string(msg.Data) + " from  NATS server ID - " + conn.ConnectedServerId())
	})

	runtime.Goexit()
}
