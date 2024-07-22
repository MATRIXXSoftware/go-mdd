package field

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneNo1(t *testing.T) {
	value := "*1234567890123#"
	PhoneNo, err := NewPhoneNo([]byte(value))
	assert.Nil(t, err)
	assert.Equal(t, value, string(PhoneNo))
}

func TestInvalidPhoneNo1(t *testing.T) {
	value := "aaa"
	_, err := NewPhoneNo([]byte(value))
	assert.Equal(t, errors.New("bad format: 'aaa' is not a valid phone number"), err)
}
