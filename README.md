# Matrixx Domain Data

## Installation
To install the MDD library, add the following import statement to your Go code:

```go
import "github.com/matrixxsoftware/go-mdd"
```

Then, run the following command to fetch the library:
```bash
go get -u github.com/matrixxsoftware/go-mdd@latest
```

# Example

## Decode
```go
	codec := cmdc.NewCodec()

	data := "<1,18,0,-6,5222,2>[1,-20,(5:value)]"
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	decoded.Containers[0].GetField(0).Type = field.UInt8
	decoded.Containers[0].GetField(1).Type = field.UInt32
	decoded.Containers[0].GetField(2).Type = field.Int32
	decoded.Containers[0].GetField(3).Type = field.String

	v, err := decoded.Containers[0].GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, uint8(1), v)

	v, err = decoded.Containers[0].GetField(1).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, nil, v)

	v, err = decoded.Containers[0].GetField(2).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(-20), v)

	v, err = decoded.Containers[0].GetField(3).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, "value", v)
```

## Encode
```go
	codec := cmdc.NewCodec()

	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{Version: 1, TotalField: 18, Depth: 0, Key: -6, SchemaVersion: 5222, ExtVersion: 2},
				Fields: []mdd.Field{
					*mdd.NewBasicField(uint8(1)),
					*mdd.NewNullField(field.UInt32),
					*mdd.NewBasicField(int32(-20)),
					*mdd.NewBasicField("value"),
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,,-20,(5:value)]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
```

## Server
```go
	codec := cmdc.NewCodec()

	transport, err := tcp.NewServerTransport("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	server := &mdd.Server{
		Codec:     codec,
		Transport: transport,
	}

	server.MessageHandler(func(containers *mdd.Containers) *mdd.Containers {
        ... 
		repsonse := mdd.Containers{}
        return response
	})

	err := transport.Listen()
	if err != nil {
		panic(err)
	}

```


## Client
```go
	codec := cmdc.NewCodec()

	transport, err := tcp.NewClientTransport("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	client := &mdd.Client{
		Codec:     codec,
		Transport: transport,
	}

	request := mdd.Containers{}

	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}
```