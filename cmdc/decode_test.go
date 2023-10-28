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
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "abc", container.GetField(1).String())
	assert.Equal(t, "foo", container.GetField(2).String())
	assert.Equal(t, "bar", container.GetField(3).String())
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
	assert.Equal(t, "", container.GetField(0).String())
	assert.Equal(t, "(6:value2)", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "2021-09-07T08:00:25.000001Z", container.GetField(3).String())
	assert.Equal(t, "2021-10-31", container.GetField(4).String())
	assert.Equal(t, "09:13:02.667997Z", container.GetField(5).String())
	assert.Equal(t, "88", container.GetField(6).String())
	assert.Equal(t, "5.5", container.GetField(7).String())
	assert.Equal(t, "", container.GetField(8).String())
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
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "abc", container0.GetField(1).String())
	assert.Equal(t, "foo", container0.GetField(2).String())
	assert.Equal(t, "bar", container0.GetField(3).String())

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
	assert.Equal(t, "2", container1.GetField(0).String())
	assert.Equal(t, "def", container1.GetField(1).String())
}

func TestDecodeNestedContainers(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,abc,<1,2,0,452,5222,2>[100],bar]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "abc", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100]", container0.GetField(2).String())
	assert.Equal(t, "bar", container0.GetField(3).String())
}

func TestDecodeFieldWithReservedCharacter(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,2,(10:v[<ue(obar),4,,6]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2", container0.GetField(1).String())
	assert.Equal(t, "(10:v[<ue(obar)", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())
	assert.Equal(t, "", container0.GetField(4).String())
	assert.Equal(t, "6", container0.GetField(5).String())
}

func TestInvalidHeader(t *testing.T) {
	mdc := "<1,18,-6,5222,2>[1,abc,foo,bar]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header, 6 fields expected"), err)
}
