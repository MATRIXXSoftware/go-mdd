package cmdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func TestIntValue(t *testing.T) {
	f := mdd.Field{Data: []byte("107")}
	assert.Equal(t, 107, f.IntValue())
}

// func TestStringValue(t *testing.T) {
// 	f := mdd.Field{Data: []byte("(6:foobar)")}
// 	assert.Equal(t, "foobar", f.StringValue())
// }
