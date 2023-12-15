package cmdc

import (
	"errors"
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

type Value struct {
	Data []byte
	V    interface{}
}

func (v Value) String() (string, error) {
	if v.V != nil {
		return v.V.(string), nil
	}
	if len(v.Data) == 0 {
		return "", nil
	}
	if v.Data[0] != '(' {
		return "", errors.New("Invalid string value")
	}
	for idx := 1; idx < len(v.Data); idx++ {
		c := v.Data[idx]
		if c == ':' {
			temp := v.Data[1:idx]
			len, err := bytesToInt(temp)
			if err != nil {
				panic("Invalid string length")
			}
			v.V = string(v.Data[idx+1 : idx+1+len])
			return v.V.(string), nil
		}
	}
	return "", errors.New("Invalid string value")
}

func (v Value) Int32() (int32, error) {
	if v.V == nil {
		value, err := strconv.Atoi(string(v.Data))
		if err != nil {
			return 0, err
		}
		v.V = int32(value)
	}
	return v.V.(int32), nil
}

func (v Value) Float32() (float32, error) {
	if v.V == nil {
		value, err := strconv.ParseFloat(string(v.Data), 32)
		if err != nil {
			return 0, err
		}
		v.V = float32(value)
	}
	return v.V.(float32), nil
}

func (v Value) Struct() (*mdd.Containers, error) {
	if v.V == nil {
		containers, err := Decode([]byte(v.Data))
		if err != nil {
			return nil, err
		}
		v.V = containers
	}
	return v.V.(*mdd.Containers), nil
}
