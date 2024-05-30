package test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd"

	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/matrixxsoftware/go-mdd/transport/http"
	"github.com/matrixxsoftware/go-mdd/transport/tcp"

	log "github.com/sirupsen/logrus"
)

func TestTransport(t *testing.T) {
	formatter := &log.TextFormatter{}
	formatter.DisableQuote = true
	log.SetFormatter(formatter)
	log.SetLevel(log.TraceLevel)

	dict := dictionary.New()
	dict.Add(&dictionary.ContainerDefinition{
		Key:           101,
		Name:          "Request",
		SchemaVersion: 5222,
		ExtVersion:    2,
		Fields: []dictionary.FieldDefinition{
			{Name: "Field1", Type: field.Int32},
			{Name: "Field2", Type: field.String},
			{Name: "Field3", Type: field.Decimal},
			{Name: "Field4", Type: field.UInt32},
			{Name: "Field5", Type: field.UInt16},
			{Name: "Field6", Type: field.UInt64},
		},
	})
	dict.Add(&dictionary.ContainerDefinition{
		Key:           88,
		Name:          "Response",
		SchemaVersion: 5222,
		ExtVersion:    2,
		Fields: []dictionary.FieldDefinition{
			{Name: "ResultCode", Type: field.UInt32},
			{Name: "ResultMessage", Type: field.String},
		},
	})

	codec := cmdc.NewCodecWithDict(dict)

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

			server.MessageHandler(func(request *mdd.Containers) (*mdd.Containers, error) {
				// t.Logf("Server received request:\n%s", request.Dump())

				container0 := request.GetContainer(101)
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
			// t.Logf("Client received response:\n%s", response.Dump())

			container0 := response.GetContainer(88)

			// Result Code
			v, err := container0.GetField(0).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, uint32(0), v.(uint32))

			// Result Message
			v, err = container0.GetField(1).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, "OK", v.(string))

			// Does not exist
			v, err = container0.GetField(2).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, nil, v)
		})
	}

}
