package main

import (
	"math/rand"
	"time"

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

	dict := loadDictionary()
	codec := cmdc.NewCodecWithDict(dict)

	addr := "0.0.0.0:14060"
	log.Infof("Server listening on %s", addr)

	transport, err := tcp.NewServerTransport(addr, codec)
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	server := &mdd.Server{
		Transport: transport,
	}

	server.MessageHandler(func(request *mdd.Containers) (*mdd.Containers, error) {
		log.Infof("Server received request:\n%s", request.Dump())

		hopId := request.GetContainer(93).GetField(14).Data

		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)

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
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: []byte("")},
						{Data: hopId},
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

	err = transport.Listen()
	if err != nil {
		panic(err)
	}
}

func loadDictionary() *dictionary.Dictionary {
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

	return dict
}
