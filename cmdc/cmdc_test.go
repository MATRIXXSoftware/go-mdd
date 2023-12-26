package cmdc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

var codec = NewCodec()

func BenchmarkDecode(b *testing.B) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4]"
	for i := 0; i < b.N; i++ {
		codec.Decode([]byte(mdc))
	}
}

func BenchmarkEncode(b *testing.B) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4]"
	containers, _ := codec.Decode([]byte(mdc))
	for i := 0; i < b.N; i++ {
		codec.Encode(containers)
	}
}

func TestEncodeDecode(t *testing.T) {

	data := "<1,18,0,-6,5222,2>[1,-20,<1,2,0,452,5222,2>[(14:abcdefghijklmn),,100,5.8888,<1,3,0,400,5222,2>[]],,400000000]"
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	// Print decoded decoded
	t.Logf("decoded: %s", decoded.Dump())

	// Assert decoded
	assert.Equal(t, "1", decoded.Containers[0].GetField(0).String())
	assert.Equal(t, "-20", decoded.Containers[0].GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[(14:abcdefghijklmn),,100,5.8888,<1,3,0,400,5222,2>[]]", decoded.Containers[0].GetField(2).String())
	assert.Equal(t, "", decoded.Containers[0].GetField(3).String())
	assert.Equal(t, "400000000", decoded.Containers[0].GetField(4).String())

	// Mark field types
	decoded.Containers[0].GetField(0).Type = field.UInt8
	decoded.Containers[0].GetField(1).Type = field.Int32
	decoded.Containers[0].GetField(2).Type = field.Struct
	decoded.Containers[0].GetField(3).Type = field.UInt32
	decoded.Containers[0].GetField(4).Type = field.UInt64

	// Retrieve field 0 as uint8
	field0, err := decoded.Containers[0].GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, uint8(1), field0)

	// Retrieve field 1 as int32
	field1, err := decoded.Containers[0].GetField(1).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(-20), field1)

	// Retrieve field 2 as struct
	field2, err := decoded.Containers[0].GetField(2).GetValue()
	assert.Nil(t, err)

	{
		nested := field2.(*mdd.Containers)
		t.Logf("nested: %s", nested.Dump())

		// Assert Nested
		assert.Equal(t, "(14:abcdefghijklmn)", nested.Containers[0].GetField(0).String())
		assert.Equal(t, "", nested.Containers[0].GetField(1).String())
		assert.Equal(t, "100", nested.Containers[0].GetField(2).String())
		assert.Equal(t, "5.8888", nested.Containers[0].GetField(3).String())
		assert.Equal(t, "<1,3,0,400,5222,2>[]", nested.Containers[0].GetField(4).String())

		// Mark field types
		nested.Containers[0].GetField(0).Type = field.String
		nested.Containers[0].GetField(1).Type = field.String
		nested.Containers[0].GetField(2).Type = field.UInt32
		nested.Containers[0].GetField(3).Type = field.Decimal
		nested.Containers[0].GetField(4).Type = field.Struct

		// Retrieve nested field 0 as string
		field0, err := nested.Containers[0].GetField(0).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, "abcdefghijklmn", field0)

		// Retrieve nested field 1 as null
		field1, err := nested.Containers[0].GetField(1).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, nil, field1)

		// Retrieve nested field 2 as UInt32
		field2, err := nested.Containers[0].GetField(2).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, uint32(100), field2)

		// Retrieve nested field 3 as Decimal
		field3, err := nested.Containers[0].GetField(3).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, "5.8888", field3.(*big.Float).Text('f', -1))

		// Retrieve nested field 4 as Struct
		field4, err := nested.Containers[0].GetField(4).GetValue()
		assert.Nil(t, err)

		{
			nested2 := field4.(*mdd.Containers)
			t.Logf("nested2: %s", nested2.Dump())

			// Assert Nested2
			assert.Equal(t, "", nested2.Containers[0].GetField(0).String())
		}
	}

	// Retrieve field 3 as uint32
	field3, err := decoded.Containers[0].GetField(3).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, nil, field3)

	// Retrieve field 4 as uint64
	field4, err := decoded.Containers[0].GetField(4).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, uint64(400000000), field4)

	// Re - encode
	encoded, err := codec.Encode(decoded)
	assert.Nil(t, err)
	assert.Equal(t, data, string(encoded))
}

func TestEncodeValues(t *testing.T) {
	nested := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{Version: 1, TotalField: 3, Depth: 0, Key: 100, SchemaVersion: 5222, ExtVersion: 2},
				Fields: []mdd.Field{
					*mdd.NewBasicField(uint32(0)),
					*mdd.NewBasicField("OK"),
					*mdd.NewNullField(field.UInt32),
				},
			},
		},
	}

	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{Version: 1, TotalField: 5, Depth: 0, Key: 200, SchemaVersion: 5222, ExtVersion: 2},
				Fields: []mdd.Field{
					*mdd.NewNullField(field.Int32),
					*mdd.NewBasicField(int32(100)),
					*mdd.NewStructField(codec, &nested),
				},
			},
		},
	}

	encoded, err := codec.Encode(&containers)
	assert.Nil(t, err)
	assert.Equal(t, "<1,5,0,200,5222,2>[,100,<1,3,0,100,5222,2>[0,(2:OK),]]", string(encoded))
}
