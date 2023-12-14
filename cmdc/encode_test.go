package cmdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

func TestEncode(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("1")},
					{Data: []byte("20")},
					{Data: []byte("(5:three)")},
					{Data: []byte("400000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,20,(5:three),400000]"
	encoded, err := Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncode2(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("1")},
					{Data: []byte("20")},
					{Data: []byte("(5:three)")},
					{Data: []byte("4000")},
				},
			},
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    20,
					Depth:         0,
					Key:           -7,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("")},
					{Data: []byte("200")},
					{Data: []byte("(7:FooBar3)")},
					{Data: []byte("")},
					{Data: []byte("100000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,20,(5:three),4000]<1,20,0,-7,5222,2>[,200,(7:FooBar3),,100000]"
	encoded, err := Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodeEmptyField(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[]"
	encoded, err := Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodeKnownType(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Type: field.Int32, Value: int32(1)},
					{Type: field.String, Value: "three"},
					{Type: field.String, Value: "富爸"},
					{Data: []byte("4000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,(5:three),(6:富爸),4000]"
	encoded, err := Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodeNested(t *testing.T) {

	subContainers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    20,
					Depth:         0,
					Key:           -7,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("")},
					{Data: []byte("200")},
					{Data: []byte("(7:FooBar3)")},
					{Data: []byte("")},
					{Data: []byte("100000")},
				},
			},
		},
	}

	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("1")},
					{Type: field.Struct, Value: &subContainers},
					{Data: []byte("(5:three)")},
					{Data: []byte("4000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,<1,20,0,-7,5222,2>[,200,(7:FooBar3),,100000],(5:three),4000]"
	encoded, err := Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}
