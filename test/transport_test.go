package test

import (
	"testing"

	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/mdd"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTransport(t *testing.T) {

	codec := cmdc.NewCodec()

	// Create Server
	server, err := mdd.NewServer("localhost:8080", codec)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	server.Handler(func(request *mdd.Containers) *mdd.Containers {
		log.Infof("Server received request : %v", request)

		container0 := request.GetContainer(101)
		assert.Equal(t, "1", container0.GetField(0).String())
		assert.Equal(t, "(3:two)", container0.GetField(1).String())
		assert.Equal(t, "3.3", container0.GetField(2).String())
		assert.Equal(t, "", container0.GetField(3).String())
		assert.Equal(t, "", container0.GetField(4).String())
		assert.Equal(t, "666", container0.GetField(5).String())

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
						{Data: []byte("(2:OK)")},
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
					{Data: []byte("(3:two)")},
					{Data: []byte("3.3")},
					{Data: []byte("")},
					{Data: []byte("")},
					{Data: []byte("666")},
				},
			},
		},
	}
	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}

	log.Infof("Client received response: %v", response)

	container0 := response.GetContainer(88)
	assert.Equal(t, "0", container0.GetField(0).String())
	assert.Equal(t, "(2:OK)", container0.GetField(1).String())
}
