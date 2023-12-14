package cmdc

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

type Value struct {
	Data []byte
	V    interface{}
}

func (v Value) Integer() int {
	if v.V != nil {
		return v.V.(int)
	}

	v.V, _ = strconv.Atoi(string(v.Data))
	return v.V.(int)
}

func (v Value) Int32() int32 {
	if v.V != nil {
		return v.V.(int32)
	}

	v.V, _ = strconv.Atoi(string(v.Data))
	return v.V.(int32)
}

func (v Value) String() string {
	if v.V != nil {
		return v.V.(string)
	}

	if len(v.Data) == 0 {
		return ""
	}
	if v.Data[0] != '(' {
		panic("Invalid string value")
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
			return v.V.(string)
		}
	}

	panic("Invalid string value")
}

func (v Value) Struct() *mdd.Containers {
	if v.V != nil {
		return v.V.(*mdd.Containers)
	}

	// TODO imeplement

	panic("Not implemented")
}
