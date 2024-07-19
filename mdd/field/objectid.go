package field

import (
	"fmt"
	"unicode"
)

type MtxObjectID string

func NewObjectID(b []byte) (MtxObjectID, error) {
	value := string(b)
	var colonCount, dashCount int
	for _, c := range value {
		if c == ':' {
			colonCount++
		} else if c == '-' {
			dashCount++
		} else if !unicode.IsDigit(c) {
			return "", fmt.Errorf("invalid ObjectId format '%s'. Parts must be numeric", value)
		}
	}
	if !(dashCount == 3 && colonCount == 0) && !(dashCount == 0 && colonCount == 3) {
		return "", fmt.Errorf("invalid ObjectId format '%s'", value)
	}
	return MtxObjectID(value), nil
}

func (o MtxObjectID) Bytes() ([]byte, error) {
	return []byte(o), nil
}
