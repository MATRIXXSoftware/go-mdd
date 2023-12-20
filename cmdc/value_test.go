package cmdc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func TestBoolTrueValue(t *testing.T) {
	v := Value{Data: []byte("1")}
	value, err := v.Bool()
	assert.Nil(t, err)
	assert.Equal(t, true, value)
}

func TestBoolFalseValue(t *testing.T) {
	v := Value{Data: []byte("0")}
	value, err := v.Bool()
	assert.Nil(t, err)
	assert.Equal(t, false, value)
}

func TestInt32Value(t *testing.T) {
	v := Value{Data: []byte("-107")}
	value, err := v.Int32()
	assert.Nil(t, err)
	assert.Equal(t, int32(-107), value)
}

func TestInt64Value(t *testing.T) {
	v := Value{Data: []byte("-107")}
	value, err := v.Int64()
	assert.Nil(t, err)
	assert.Equal(t, int64(-107), value)
}

func TestUInt32Value(t *testing.T) {
	v := Value{Data: []byte("1070")}
	value, err := v.UInt32()
	assert.Nil(t, err)
	assert.Equal(t, uint32(1070), value)
}

func TestUInt64Value(t *testing.T) {
	v := Value{Data: []byte("1070")}
	value, err := v.UInt64()
	assert.Nil(t, err)
	assert.Equal(t, uint64(1070), value)
}

func TestStringValue(t *testing.T) {
	v := Value{Data: []byte("(6:foobar)")}
	value, err := v.String()
	assert.Nil(t, err)
	assert.Equal(t, "foobar", value)
}

func TestUnicodeStringValue(t *testing.T) {
	v := Value{Data: []byte("(6:富爸)")}
	value, err := v.String()
	assert.Nil(t, err)
	assert.Equal(t, "富爸", value)
}

func TestStructValue(t *testing.T) {
	v := Value{Data: []byte("<1,10,0,235,5280,1>[1,20,300,4]")}

	containers, err := v.Struct()
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
}

func TestDecimalValue(t *testing.T) {
	v := Value{Data: []byte("3.142")}
	value, err := v.Decimal()
	assert.Nil(t, err)
	assert.Equal(t, new(big.Float).SetFloat64(3.142).Text('f', -1), value.Text('f', -1))
}

func TestDateTimeValue(t *testing.T) {
	v := Value{Data: []byte("2017-01-01T12:00:00Z")}
	value, err := v.DateTime()
	assert.Nil(t, err)
	assert.Equal(t, "2017-01-01 12:00:00 +0000 UTC", value.String())
}
