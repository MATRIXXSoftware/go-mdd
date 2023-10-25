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
				Fields: []Field{{Value: "1"}, {Value: "two"}, {Value: "three"}, {Value: "4"}},
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
				Fields: []Field{{Value: "1"}, {Value: "abc"}, {Value: "foo"}, {Value: "bar"}, {Value: "5"}, {Value: "6"}},
			},
		},
	}

	container1 := mdc.GetContainer(101)
	assert.Equal(t, mdc.Containers[0], *container1)

	container2 := mdc.GetContainer(102)
	assert.Equal(t, mdc.Containers[1], *container2)

	container3 := mdc.GetContainer(1000)
	assert.Nil(t, container3)
}
