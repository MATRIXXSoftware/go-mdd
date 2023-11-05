package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/transport/http"
)

func TestTransport(t *testing.T) {

	codec := cmdc.NewCodec()

	// Create Server
	serverTransport, err := http.NewServerTransport("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer serverTransport.Close()

	server := &mdd.Server{
		Codec:     codec,
		Transport: serverTransport,
	}

	server.Handler(func(request *mdd.Containers) *mdd.Containers {
		t.Logf("Server received request:\n%s", request.Dump())

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
		err := serverTransport.Listen()
		if err != nil {
			panic(err)
		}
	}()

	// Create Client
	clientTransport, err := http.NewClientTransport("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer clientTransport.Close()

	client := &mdd.Client{
		Codec:     codec,
		Transport: clientTransport,
	}

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

	t.Logf("Client received response:\n%s", response.Dump())

	container0 := response.GetContainer(88)
	assert.Equal(t, "0", container0.GetField(0).String())
	assert.Equal(t, "(2:OK)", container0.GetField(1).String())
}
