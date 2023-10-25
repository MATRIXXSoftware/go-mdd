package mdd

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestTransport(t *testing.T) {
	// Create Server
	server, err := NewServer("localhost:8080")
	if err != nil {
		panic(err)
	}

	go func() {
		err := server.Listen()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Create Client
	client, err := NewClient("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Send Message
	request := ExampleMessage{Field1: 10, Field2: 20}
	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}

	log.Infof("Response: %+v", response)
}
