package transport

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestTransport(t *testing.T) {
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

	client, err := NewClient("localhost:8080")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	request := ExampleMessage{Field1: 10, Field2: 20}
	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}

	log.Infof("Response: %+v", response)

	time.Sleep(100 * time.Millisecond)
}
