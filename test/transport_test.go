package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd"

	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/matrixxsoftware/go-mdd/transport/client"
	"github.com/matrixxsoftware/go-mdd/transport/protocol/http"
	"github.com/matrixxsoftware/go-mdd/transport/protocol/tcp"
	"github.com/matrixxsoftware/go-mdd/transport/server"

	log "github.com/sirupsen/logrus"
)

func TestTransport(t *testing.T) {
	formatter := &log.TextFormatter{}
	formatter.DisableQuote = true
	log.SetFormatter(formatter)
	log.SetLevel(log.TraceLevel)

	dict := dictionary.New()
	dict.Add(&dictionary.ContainerDefinition{
		Key:           93,
		Name:          "MtxMsg",
		SchemaVersion: 5222,
		ExtVersion:    1,
		Fields: []dictionary.FieldDefinition{
			{Name: "Field1", Type: field.Int32},
			{Name: "Field2", Type: field.String},
			{Name: "Field3", Type: field.Decimal},
			{Name: "Field4", Type: field.UInt32},
			{Name: "Field5", Type: field.UInt16},
			{Name: "Field6", Type: field.UInt64},
			{Name: "Field7", Type: field.UInt32},
			{Name: "Field8", Type: field.UInt32},
			{Name: "Field9", Type: field.UInt32},
			{Name: "Field10", Type: field.UInt32},
			{Name: "Field11", Type: field.UInt32},
			{Name: "Field12", Type: field.UInt32},
			{Name: "Field13", Type: field.UInt32},
			{Name: "Field14", Type: field.UInt32},
			{Name: "Field15", Type: field.UInt32},
		},
	})
	dict.Add(&dictionary.ContainerDefinition{
		Key:           235,
		Name:          "MtxRequest",
		SchemaVersion: 5222,
		ExtVersion:    1,
		Fields: []dictionary.FieldDefinition{
			{Name: "Version", Type: field.UInt16},
			{Name: "EventTime", Type: field.DateTime},
		},
	})
	dict.Add(&dictionary.ContainerDefinition{
		Key:           236,
		Name:          "MtxResponse",
		SchemaVersion: 5222,
		ExtVersion:    1,
		Fields: []dictionary.FieldDefinition{
			{Name: "RouteId", Type: field.UInt16},
			{Name: "Result", Type: field.UInt32},
			{Name: "ResultText", Type: field.String},
		},
	})

	codec := cmdc.NewCodecWithDict(dict)

	transports := []struct {
		name               string
		addr               string
		newServerTransport func(string) (server.Transport, error)
		newClientTransport func(string) (client.Transport, error)
	}{
		{
			"TCP",
			"localhost:8080",
			func(addr string) (server.Transport, error) {
				return tcp.NewServerTransport(addr, codec)
			},
			func(addr string) (client.Transport, error) {
				return tcp.NewClientTransport(addr, codec)
			},
		},
		{
			"TCP+TLS",
			"localhost:8081",
			func(addr string) (server.Transport, error) {
				return tcp.NewServerTransport(addr, codec,
					server.WithTls(server.TLS{
						Enabled:        true,
						SelfSignedCert: true,
					}),
				)

			},
			func(addr string) (client.Transport, error) {
				return tcp.NewClientTransport(addr, codec,
					client.WithTls(client.TLS{
						Enabled:            true,
						InsecureSkipVerify: true,
					}),
				)
			},
		},
		{
			"HTTP",
			"localhost:8082",
			func(addr string) (server.Transport, error) {
				return http.NewServerTransport(addr, codec)
			},
			func(addr string) (client.Transport, error) {
				return http.NewClientTransport(addr, codec)
			},
		},
		{
			"HTTP+TLS",
			"localhost:8083",
			func(addr string) (server.Transport, error) {
				return http.NewServerTransport(addr, codec,
					server.WithTls(server.TLS{
						Enabled:        true,
						SelfSignedCert: true,
					}),
				)
			},
			func(addr string) (client.Transport, error) {
				return http.NewClientTransport(addr, codec,
					client.WithTls(client.TLS{
						Enabled:            true,
						InsecureSkipVerify: true,
					}),
				)
			},
		},
	}

	for _, tt := range transports {
		t.Run(tt.name, func(t *testing.T) {
			serverTransport, err := tt.newServerTransport(tt.addr)
			if err != nil {
				t.Fatalf("failed to create server transport: %v", err)
			}
			defer serverTransport.Close()

			server := &server.Server{
				Transport: serverTransport,
			}

			server.MessageHandler(func(request *mdd.Containers) (*mdd.Containers, error) {
				t.Logf("Server received request:\n%s", request.Dump())

				container0 := request.GetContainer(93)
				// Field 1
				v, err := container0.GetField(0).GetValue()
				assert.Nil(t, err)
				assert.Equal(t, int32(1), v.(int32))

				// Field 2
				v, err = container0.GetField(1).GetValue()
				assert.Nil(t, err)
				assert.Equal(t, "two", v.(string))

				// Field 3
				v, err = container0.GetField(2).GetValue()
				assert.Nil(t, err)
				assert.Equal(t, "3.3", v.(*big.Float).Text('f', 1))

				// Field 4
				v, err = container0.GetField(3).GetValue()
				assert.Nil(t, err)
				assert.Equal(t, nil, v)

				// Field 5
				v, err = container0.GetField(4).GetValue()
				assert.Nil(t, err)
				assert.Equal(t, nil, v)

				// Field 6
				v, err = container0.GetField(5).GetValue()
				assert.Nil(t, err)
				assert.Equal(t, uint64(666), v.(uint64))

				return &mdd.Containers{
					Containers: []mdd.Container{
						{
							Header: mdd.Header{
								Version:       1,
								TotalField:    14,
								Depth:         0,
								Key:           93,
								SchemaVersion: 5222,
								ExtVersion:    1,
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
						{
							Header: mdd.Header{
								Version:       1,
								TotalField:    3,
								Depth:         0,
								Key:           236,
								SchemaVersion: 5222,
								ExtVersion:    1,
							},
							Fields: []mdd.Field{
								{Data: []byte("1")},
								{Data: []byte("0")},
								{Data: []byte("(2:OK)")},
							},
						},
					},
				}, nil
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
			clientTransport, err := tt.newClientTransport(tt.addr)
			if err != nil {
				panic(err)
			}
			defer clientTransport.Close()

			client := &client.Client{
				Transport: clientTransport,
			}

			// Send Message
			request := mdd.Containers{
				Containers: []mdd.Container{
					{
						Header: mdd.Header{
							Version:       1,
							TotalField:    14,
							Depth:         0,
							Key:           93,
							SchemaVersion: 5222,
							ExtVersion:    1,
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
					{
						Header: mdd.Header{
							Version:       1,
							TotalField:    2,
							Depth:         0,
							Key:           235,
							SchemaVersion: 5222,
							ExtVersion:    1,
						},
						Fields: []mdd.Field{
							{Data: []byte("1")},
							{Data: []byte("2021-01-01T00:00:00Z")},
						},
					},
				},
			}
			response, err := client.SendMessage(context.Background(), &request)
			if err != nil {
				panic(err)
			}
			t.Logf("Client received response:\n%s", response.Dump())

			container0 := response.GetContainer(236)

			// RouteId
			v, err := container0.GetField(0).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, uint16(1), v.(uint16))

			// Result Code
			v, err = container0.GetField(1).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, uint32(0), v.(uint32))

			// Result Message
			v, err = container0.GetField(2).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, "OK", v.(string))

			// Does not exist
			v, err = container0.GetField(3).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, nil, v)
		})
	}
}
