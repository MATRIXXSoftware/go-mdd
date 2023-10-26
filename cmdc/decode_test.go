package cmdc

import (
	"errors"
	"fmt"
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

	expectedFields := []mdd.Field{
		{Data: []byte("1")},
		{Data: []byte("abc")},
		{Data: []byte("foo")},
		{Data: []byte("bar")},
	}
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
		{Data: []byte("")},
		{Data: []byte("(6:value2)")},
		{Data: []byte("3")},
		{Data: []byte("2021-09-07T08:00:25.000001Z")},
		{Data: []byte("2021-10-31")},
		{Data: []byte("09:13:02.667997Z")},
		{Data: []byte("88")},
		{Data: []byte("5.5")},
		{Data: []byte("")}}
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
	expectedFields := []mdd.Field{
		{Data: []byte("1")},
		{Data: []byte("abc")},
		{Data: []byte("foo")},
		{Data: []byte("bar")},
	}
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
	expectedFields = []mdd.Field{
		{Data: []byte("2")},
		{Data: []byte("def")},
	}
	assert.Equal(t, expectedFields, container1.Fields)
}

func TestDecodeNestedContainers(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,<1,2,0,452,5222,2>[100],bar]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	fmt.Printf("container0: %v\n", container0)

	expectedFields := []mdd.Field{
		{Data: []byte("1")},
		{Data: []byte("abc")},
		{Data: []byte("<1,2,0,452,5222,2>[100]")},
		{Data: []byte("bar")},
	}
	assert.Equal(t, expectedFields, container0.Fields)
}

func TestInvalidHeader(t *testing.T) {
	mdc := "<1,18,-6,5222,2>[1,abc,foo,bar]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header, 6 fields expected"), err)
}
