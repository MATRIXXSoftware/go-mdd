package cmdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func TestIntValue(t *testing.T) {
	v := Value{Data: []byte("107")}
	assert.Equal(t, 107, v.Integer())
}

func TestStringValue(t *testing.T) {
	v := Value{Data: []byte("(6:foobar)")}
	assert.Equal(t, "foobar", v.String())
}

func extractString(f *mdd.Field) string {
	return f.String()
}
