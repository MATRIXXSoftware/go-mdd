package field

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectID1(t *testing.T) {
	value := "1-5-1-1"
	objectID, err := NewObjectID(value)
	assert.Nil(t, err)
	assert.Equal(t, value, string(objectID))
}

func TestObjectID2(t *testing.T) {
	value := "1:5:1:1"
	objectID, err := NewObjectID(value)
	assert.Nil(t, err)
	assert.Equal(t, value, string(objectID))
}

func TestInvalidObjectID1(t *testing.T) {
	value := "15-1-1"
	_, err := NewObjectID(value)
	assert.Equal(t, errors.New("invalid ObjectId format '15-1-1'"), err)
}

func TestInvalidObjectID2(t *testing.T) {
	value := "1-5-1-1-1"
	_, err := NewObjectID(value)
	assert.Equal(t, errors.New("invalid ObjectId format '1-5-1-1-1'"), err)
}
