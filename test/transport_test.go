package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/transport/http"
	"github.com/matrixxsoftware/go-mdd/transport/tcp"
)

func TestTransport(t *testing.T) {

	codec := cmdc.NewCodec()

	transports := []struct {
		name               string
		newServerTransport func(string) (mdd.ServerTransport, error)
		newClientTransport func(string) (mdd.ClientTransport, error)
	}{
		{
			"TCP",
			func(addr string) (mdd.ServerTransport, error) {
				return tcp.NewServerTransport(addr)
			},
			func(addr string) (mdd.ClientTransport, error) {
				return tcp.NewClientTransport(addr)
			},
		},
		{
			"HTTP",
			func(addr string) (mdd.ServerTransport, error) {
				return http.NewServerTransport(addr)
			},
			func(addr string) (mdd.ClientTransport, error) {
				return http.NewClientTransport(addr)
			},
		},
	}

	for _, tt := range transports {
		t.Run(tt.name, func(t *testing.T) {
			serverTransport, err := tt.newServerTransport("localhost:8080")
			if err != nil {
				t.Fatalf("failed to create server transport: %v", err)
			}
			defer serverTransport.Close()

			server := &mdd.Server{
				Codec:     codec,
				Transport: serverTransport,
			}

			server.MessageHandler(func(request *mdd.Containers) *mdd.Containers {
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

			// Add a small delay for server to start
			time.Sleep(100 * time.Millisecond)

			// Create Client
			clientTransport, err := tt.newClientTransport("localhost:8080")
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

			assert.Equal(t, 0, container0.GetField(0).Value.Integer())
			assert.Equal(t, "OK", container0.GetField(1).Value.String())
		})
	}

}
