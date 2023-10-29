package mdd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainer(t *testing.T) {

	mdc := Containers{
		Containers: []Container{
			{
				Header: Header{
					Version:       1,
					TotalField:    5,
					Depth:         0,
					Key:           101,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []Field{
					{Data: []byte("1")},
					{Data: []byte("2.2")},
					{Data: []byte("(5:three)")},
					{Data: []byte("40000")},
				},
			},
			{
				Header: Header{
					Version:       1,
					TotalField:    6,
					Depth:         0,
					Key:           102,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []Field{
					{Data: []byte("")},
					{Data: []byte("(3:abc)")},
					{Data: []byte("(3:foo)")},
					{Data: []byte("(3:bar)")},
					{Data: []byte("5.0")},
					{Data: []byte("")},
					{Data: []byte("777")},
				},
			},
		},
	}

	container0 := mdc.GetContainer(101)
	assert.Equal(t, mdc.Containers[0], *container0)
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2.2", container0.GetField(1).String())
	assert.Equal(t, "(5:three)", container0.GetField(2).String())
	assert.Equal(t, "40000", container0.GetField(3).String())

	container1 := mdc.GetContainer(102)
	assert.Equal(t, mdc.Containers[1], *container1)
	assert.Equal(t, "", container1.GetField(0).String())
	assert.Equal(t, "(3:abc)", container1.GetField(1).String())
	assert.Equal(t, "(3:foo)", container1.GetField(2).String())
	assert.Equal(t, "(3:bar)", container1.GetField(3).String())
	assert.Equal(t, "5.0", container1.GetField(4).String())
	assert.Equal(t, "", container1.GetField(5).String())
	assert.Equal(t, "777", container1.GetField(6).String())

	// Field 6 does not exist
	assert.Nil(t, container1.GetField(7))

	// Container 2 does not exist
	container2 := mdc.GetContainer(1000)
	assert.Nil(t, container2)
}
