package test

import (
	"testing"

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
	defer server.Close()

	server.Handler(func(containers *mdd.Containers) *mdd.Containers {
		log.Infof("Server received request : %v", containers)
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
					Fields: []mdd.Field{
						{Data: []byte("0")},
						{Data: []byte("OK")},
					},
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
				Fields: []mdd.Field{
					{Data: []byte("1")},
					{Data: []byte("two")},
					{Data: []byte("three")},
					{Data: []byte("4")},
				},
			},
		},
	}
	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}

	log.Infof("Client received response: %v", response)
}
