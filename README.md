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

# Codec Example

## Decode
```go
	codec := cmdc.NewCodec()

	data := "<1,18,0,-6,5222,2>[1,,-20,(5:value),{10,20}]"
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	decoded.Containers[0].LoadDefinition(&dictionary.ContainerDefinition{
		Fields: []dictionary.FieldDefinition{
			{Type: field.UInt8},
			{Type: field.UInt32},
			{Type: field.Int32},
			{Type: field.String},
			{Type: field.Int32, IsMulti: true},
		},
	})

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

	v, err = decoded.Containers[0].GetField(4).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, []int32{10, 20}, v)
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
					*mdd.NewBasicListField([]int32{10, 20}),
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,,-20,(5:value),{10,20}]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
```

# Transport Example
- [Server Example](examples/server/server.go)
- [Client Example](examples/client/client.go)