package mdd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple1(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,foo,bar]"
	container, err := Decode([]byte(mdc))
	assert.Nil(t, err)

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
		"2021-10-31,09:13:02.667997Z,88,5.5,]"
	container, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	expectedHeader := Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container.Header)

	expectedFields := []Field{
		{""},
		{"(6:value2)"},
		{"3"},
		{"2021-09-07T08:00:25.000001Z"},
		{"2021-10-31"},
		{"09:13:02.667997Z"},
		{"88"},
		{"5.5"},
		{""}}
	assert.Equal(t, expectedFields, container.Fields)
}

func TestSimple3(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,foo,bar]<1,5,0,-7,5222,2>[2,def]"
	container, err := Decode([]byte(mdc))

	assert.Nil(t, err)

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

func TestInvalidHeader(t *testing.T) {
	mdc := "<1,18,-6,5222,2>[1,abc,foo,bar]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header"), err)
}

// func TestEncode(t *testing.T) {
// 	container := Container{
// 		Header: Header{
// 			Version:       1,
// 			TotalField:    18,
// 			Depth:         0,
// 			Key:           -6,
// 			SchemaVersion: 5222,
// 			ExtVersion:    2,
// 		},
// 		Fields: []Field{{"1"}, {"abc"}, {"foo"}, {"bar"}},
// 	}

// 	expected := "<1,18,0,-6,5222,2>[1,abc,foo,bar]"
// 	encoded, err := Encode(container)
// 	assert.Nil(t, err)
// 	assert.Equal(t, expected, encoded)
// }
