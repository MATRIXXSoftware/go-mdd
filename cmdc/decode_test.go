package cmdc

import (
	"errors"
	"testing"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/stretchr/testify/assert"
)

func TestDecodeSingleContainer1(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,foo,bar]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container.Header)

	expectedFields := []mdd.Field{{Value: "1"}, {Value: "abc"}, {Value: "foo"}, {Value: "bar"}}
	assert.Equal(t, expectedFields, container.Fields)
}

func TestDecodeSingleContainer2(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[,(6:value2),3,2021-09-07T08:00:25.000001Z," +
		"2021-10-31,09:13:02.667997Z,88,5.5,]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container.Header)

	expectedFields := []mdd.Field{
		{Value: ""},
		{Value: "(6:value2)"},
		{Value: "3"},
		{Value: "2021-09-07T08:00:25.000001Z"},
		{Value: "2021-10-31"},
		{Value: "09:13:02.667997Z"},
		{Value: "88"},
		{Value: "5.5"},
		{Value: ""}}
	assert.Equal(t, expectedFields, container.Fields)
}

func TestDecodeContainers(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,foo,bar]<1,5,0,-7,5222,2>[2,def]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container0.Header)
	expectedFields := []mdd.Field{{Value: "1"}, {Value: "abc"}, {Value: "foo"}, {Value: "bar"}}
	assert.Equal(t, expectedFields, container0.Fields)

	container1 := containers.Containers[1]

	expectedHeader = mdd.Header{
		Version:       1,
		TotalField:    5,
		Depth:         0,
		Key:           -7,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container1.Header)
	expectedFields = []mdd.Field{{Value: "2"}, {Value: "def"}}
	assert.Equal(t, expectedFields, container1.Fields)
}

func TestInvalidHeader(t *testing.T) {
	mdc := "<1,18,-6,5222,2>[1,abc,foo,bar]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header, 6 fields expected"), err)
}
