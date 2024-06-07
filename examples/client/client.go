package main

import (
	"github.com/matrixxsoftware/go-mdd/cmdc"
	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/matrixxsoftware/go-mdd/transport/tcp"

	log "github.com/sirupsen/logrus"
)

func main() {
	formatter := &log.TextFormatter{
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
	}
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

	addr := "0.0.0.0:14060"
	log.Infof("Server listening on %s", addr)

	transport, err := tcp.NewClientTransport(addr, codec)
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	client := &mdd.Client{
		Codec:     codec,
		Transport: transport,
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

	// for {
	response, err := client.SendMessage(&request)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Client received response:\n%s", response.Dump())
	//
	// 	<-time.After(2 * time.Second)
	// }
}
