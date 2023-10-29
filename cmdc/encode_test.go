package cmdc

import (
	"fmt"
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

	fmt.Printf("encoded = %s\n", encoded)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}
