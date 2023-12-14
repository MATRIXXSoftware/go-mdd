package cmdc

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

type Value struct {
	Field *mdd.Field
}

func (v *Value) Integer() int {
	f := v.Field
	if f.Value != nil {
		return f.Value.(int)
	}

	f.Value, _ = strconv.Atoi(f.String())
	return f.Value.(int)
}

func (v *Value) String() string {
	f := v.Field
	data := f.Data

	if len(data) == 0 {
		return ""
	}
	if data[0] != '(' {
		panic("Invalid string value")
	}

	for idx := 1; idx < len(data); idx++ {
		c := data[idx]
		if c == ':' {
			temp := data[1:idx]
			len, err := bytesToInt(temp)
			if err != nil {
				panic("Invalid string length")
			}
			return string(data[idx+1 : idx+1+len])
		}
	}

	panic("Invalid string value")
}
