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
	request := Containers{
		Containers: []Container{
			{
				Header: Header{
					Version:       1,
					TotalField:    5,
					Depth:         0,
					Key:           101,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []Field{{Value: "1"}, {Value: "two"}, {Value: "three"}, {Value: "4"}},
			},
		},
	}
	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}

	log.Infof("Response: %+v", response)
}
