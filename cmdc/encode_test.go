package cmdc

import (
	"testing"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	container := mdd.Container{
		Header: mdd.Header{
			Version:       1,
			TotalField:    18,
			Depth:         0,
			Key:           -6,
			SchemaVersion: 5222,
			ExtVersion:    2,
		},
		Fields: []mdd.Field{{Value: "1"}, {Value: "abc"}, {Value: "foo"}, {Value: "bar"}},
	}

	expected := "<1,18,0,-6,5222,2>[1,abc,foo,bar]"
	encoded, err := Encode(container)
	assert.Nil(t, err)
	assert.Equal(t, expected, encoded)
}
