package cmdc

import (
	"testing"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/stretchr/testify/assert"
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
