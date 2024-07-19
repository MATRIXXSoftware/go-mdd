package field

import (
	"fmt"
)

type MtxPhoneNo string

func NewPhoneNo(b []byte) (MtxPhoneNo, error) {
	value := string(b)
	if len(value) > 15 {
		return "", fmt.Errorf("phone number too long. Value limited to 15 digits")
	}
	for _, c := range value {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '*', '#':
			// valid character, do nothing
		default:
			return "", fmt.Errorf("bad format: '%s' is not a valid phone number", value)
		}
	}
	return MtxPhoneNo(value), nil
}

func (o MtxPhoneNo) Bytes() ([]byte, error) {
	return []byte(o), nil
}
