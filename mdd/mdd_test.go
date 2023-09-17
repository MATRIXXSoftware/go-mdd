package mdd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple1(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,foo,bar]"
	encode()
	container := decode(mdc)

	expectedHeader := Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container.Header)

	expectedFields := []Field{{"1"}, {"abc"}, {"foo"}, {"bar"}}
	assert.Equal(t, expectedFields, container.Fields)
}

func TestSimple2(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[,(6:value2),3,2021-09-07T08:00:25.000001Z," +
		"2021-10-31,09:13:02.667997Z,88,5.5,"
	encode()
	container := decode(mdc)
	fmt.Println("Decoded Container:")
	fmt.Printf("Header: %+v\n", container.Header)
	fmt.Printf("Fields: %+v\n", container.Fields)
}
