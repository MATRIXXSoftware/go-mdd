package cmdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func TestIntValue(t *testing.T) {
	f := mdd.Field{Data: []byte("107")}
	v := Value{Field: &f}
	assert.Equal(t, 107, v.Integer())
}

func TestStringValue(t *testing.T) {
	f := mdd.Field{Data: []byte("(6:foobar)")}
	v := Value{Field: &f}
	assert.Equal(t, "foobar", v.String())
}

func extractString(f *mdd.Field) string {
	return f.String()
}
