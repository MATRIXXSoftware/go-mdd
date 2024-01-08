package cmdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func TestBoolTrueValue(t *testing.T) {
	data := []byte("1")
	value, err := decodeBoolValue(data)
	assert.Nil(t, err)
	assert.Equal(t, true, value)

	encoded, err := encodeBoolValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestBoolFalseValue(t *testing.T) {
	data := []byte("0")
	value, err := decodeBoolValue(data)
	assert.Nil(t, err)
	assert.Equal(t, false, value)

	encoded, err := encodeBoolValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestInt8Value(t *testing.T) {
	data := []byte("-125")
	value, err := decodeInt8Value(data)
	assert.Nil(t, err)
	assert.Equal(t, int8(-125), value)

	encoded, err := encodeInt8Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestInt16Value(t *testing.T) {
	data := []byte("1070")
	value, err := decodeInt16Value(data)
	assert.Nil(t, err)
	assert.Equal(t, int16(1070), value)

	encoded, err := encodeInt16Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestInt32Value(t *testing.T) {
	data := []byte("-107")
	value, err := decodeInt32Value(data)
	assert.Nil(t, err)
	assert.Equal(t, int32(-107), value)

	encoded, err := encodeInt32Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestInt32ListValue(t *testing.T) {
	data := []byte("{1,2,3,4,5}")
	list, err := decodeListValue(data, decodeInt32Value)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(list))
	assert.Equal(t, int32(1), list[0])
	assert.Equal(t, int32(2), list[1])
	assert.Equal(t, int32(3), list[2])
	assert.Equal(t, int32(4), list[3])
	assert.Equal(t, int32(5), list[4])

	encoded, err := encodeListValue(list, encodeInt32Value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestInt64Value(t *testing.T) {
	data := []byte("81345123666616")
	value, err := decodeInt64Value(data)
	assert.Nil(t, err)
	assert.Equal(t, int64(81345123666616), value)

	encoded, err := encodeInt64Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestUInt8Value(t *testing.T) {
	data := []byte("200")
	value, err := decodeUInt8Value(data)
	assert.Nil(t, err)
	assert.Equal(t, uint8(200), value)

	encoded, err := encodeUInt8Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestUInt16Value(t *testing.T) {
	data := []byte("1070")
	value, err := decodeUInt16Value(data)
	assert.Nil(t, err)
	assert.Equal(t, uint16(1070), value)

	encoded, err := encodeUInt16Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestUInt32Value(t *testing.T) {
	data := []byte("10888888")
	value, err := decodeUInt32Value(data)
	assert.Nil(t, err)
	assert.Equal(t, uint32(10888888), value)

	encoded, err := encodeUInt32Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestUInt64Value(t *testing.T) {
	data := []byte("81345123666616")
	value, err := decodeUInt64Value(data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(81345123666616), value)

	encoded, err := encodeUInt64Value(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestStringValue(t *testing.T) {
	data := []byte("(6:foobar)")
	value, err := decodeStringValue(data)
	assert.Nil(t, err)
	assert.Equal(t, "foobar", value)

	encoded, err := encodeStringValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestUnicodeStringValue(t *testing.T) {
	data := []byte("(6:富爸)")
	value, err := decodeStringValue(data)
	assert.Nil(t, err)
	assert.Equal(t, "富爸", value)

	encoded, err := encodeStringValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestStringListValue(t *testing.T) {
	data := []byte("{(6:value1),(7:value20),(8:value300)}")
	list, err := decodeListValue(data, decodeStringValue)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "value1", list[0])
	assert.Equal(t, "value20", list[1])
	assert.Equal(t, "value300", list[2])

	encoded, err := encodeListValue(list, encodeStringValue)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestStructValue(t *testing.T) {
	data := []byte("<1,10,0,235,5280,1>[1,20,300,4]")

	containers, err := decodeStructValue(codec, data)
	assert.Nil(t, err)
	assert.NotNil(t, containers)

	container := containers.Containers[0]
	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    10,
		Depth:         0,
		Key:           235,
		SchemaVersion: 5280,
		ExtVersion:    1,
	}
	assert.Equal(t, expectedHeader, container.Header)
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "20", container.GetField(1).String())
	assert.Equal(t, "300", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())

	encoded, err := encodeStructValue(codec, containers)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestStructListValue(t *testing.T) {
	data := []byte("{<1,10,0,235,5280,1>[100],<1,10,0,235,5280,1>[200],<1,10,0,235,5280,1>[300]}")

	list, err := decodeListValue(data, func(b []byte) (*mdd.Containers, error) {
		return decodeStructValue(codec, b)
	})
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    10,
		Depth:         0,
		Key:           235,
		SchemaVersion: 5280,
		ExtVersion:    1,
	}
	assert.Equal(t, expectedHeader, list[0].Containers[0].Header)
	assert.Equal(t, "100", list[0].Containers[0].GetField(0).String())
	assert.Equal(t, "200", list[1].Containers[0].GetField(0).String())
	assert.Equal(t, "300", list[2].Containers[0].GetField(0).String())

	encoded, err := encodeListValue(list, func(v *mdd.Containers) ([]byte, error) {
		return encodeStructValue(codec, v)
	})
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestDecimalValue(t *testing.T) {
	data := []byte("3.142")
	value, err := decodeDecimalValue(data)
	assert.Nil(t, err)
	assert.Equal(t, "3.142", value.Text('f', -1))

	encoded, err := encodeDecimalValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestDateTimeValue(t *testing.T) {
	data := []byte("2024-01-08T21:54:10Z")
	value, err := decodeDateTimeValue(data)
	assert.Nil(t, err)
	assert.Equal(t, "2024-01-08 21:54:10 +0000 UTC", value.String())

	encoded, err := encodeDateTimeValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestDateValue(t *testing.T) {
	data := []byte("2024-01-08")
	value, err := decodeDateValue(data)
	assert.Nil(t, err)
	assert.Equal(t, "2024-01-08 00:00:00 +0000 UTC", value.String())

	encoded, err := encodeDateValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestTimeValue(t *testing.T) {
	data := []byte("21:54:10")
	value, err := decodeTimeValue(data)
	assert.Nil(t, err)
	assert.Equal(t, "0000-01-01 21:54:10 +0000 UTC", value.String())

	encoded, err := encodeTimeValue(value)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}
