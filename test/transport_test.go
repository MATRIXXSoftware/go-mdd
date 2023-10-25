package test

import (
	"testing"
	"time"

	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/mdd"
	log "github.com/sirupsen/logrus"
)

func TestTransport(t *testing.T) {

	codec := cmdc.NewCodec()

	// Create Server
	server, err := mdd.NewServer("localhost:8080", codec)
	if err != nil {
		panic(err)
	}

	server.Handler(func(containers *mdd.Containers) *mdd.Containers {
		log.Infof("Server received request: %+v", containers)
		return &mdd.Containers{
			Containers: []mdd.Container{
				{
					Header: mdd.Header{
						Version:       1,
						TotalField:    3,
						Depth:         0,
						Key:           88,
						SchemaVersion: 5222,
						ExtVersion:    2,
					},
					Fields: []mdd.Field{{Value: "0"}, {Value: "OK"}},
				},
			},
		}
	})

	go func() {
		err := server.Listen()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Create Client
	client, err := mdd.NewClient("localhost:8080", codec)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Send Message
	request := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    5,
					Depth:         0,
					Key:           101,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{{Value: "1"}, {Value: "two"}, {Value: "three"}, {Value: "4"}},
			},
		},
	}
	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}

	log.Infof("Client received response: %+v", response)

	time.Sleep(100 * time.Millisecond)

	// server.Close()
	// client.Close()
}
