package cmdc

import (
	"errors"
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

// String
type StringValue string

func DecodeStringValue(b []byte) (mdd.Value, error) {
	if len(b) == 0 {
		return StringValue(""), nil
	}
	if b[0] != '(' {
		return StringValue(""), errors.New("Invalid string value")
	}
	for idx := 1; idx < len(b); idx++ {
		c := b[idx]
		if c == ':' {
			temp := b[1:idx]
			len, err := bytesToInt(temp)
			if err != nil {
				panic("Invalid string length")
			}
			str := string(b[idx+1 : idx+1+len])
			return StringValue(str), nil
		}
	}
	return StringValue(""), errors.New("Invalid string value")
}

func (v StringValue) Type() field.Type {
	return field.String
}

func (v StringValue) Serialize() []byte {
	return []byte(v)
}

// Int32
type Int32Value int32

func DecodeInt32Value(b []byte) (mdd.Value, error) {
	v, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return Int32Value(0), err
	}
	return Int32Value(int32(v)), nil
}

func (v Int32Value) Type() field.Type {
	return field.Int32
}

func (v Int32Value) Serialize() []byte {
	return []byte(string(v))
}

// type CmdcFieldCodec struct{}
//
// func (c CmdcFieldCodec) Decode(f *mdd.Field) (mdd.Value, error) {
// 	switch f.Type {
// 	case field.String:
// 		return DecodeStringValue(f.Data)
// 	case field.Int32:
// 		return DecodeInt32Value(f.Data)
// 	default:
// 		return nil, errors.New("Unsupported field type")
// 	}
// }

//
// import (
// 	"errors"
// 	"math/big"
// 	"strconv"
// 	"time"
//
// 	"github.com/matrixxsoftware/go-mdd/mdd"
// )
//
// type Value struct {
// 	Data []byte
// 	V    interface{}
// }
//
// func (v Value) Bool() (bool, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseBool(string(v.Data))
// 		if err != nil {
// 			return false, err
// 		}
// 		v.V = value
// 	}
// 	return v.V.(bool), nil
// }
//
// func (v Value) String() (string, error) {
// 	if v.V != nil {
// 		return v.V.(string), nil
// 	}
// 	if len(v.Data) == 0 {
// 		return "", nil
// 	}
// 	if v.Data[0] != '(' {
// 		return "", errors.New("Invalid string value")
// 	}
// 	for idx := 1; idx < len(v.Data); idx++ {
// 		c := v.Data[idx]
// 		if c == ':' {
// 			temp := v.Data[1:idx]
// 			len, err := bytesToInt(temp)
// 			if err != nil {
// 				panic("Invalid string length")
// 			}
// 			v.V = string(v.Data[idx+1 : idx+1+len])
// 			return v.V.(string), nil
// 		}
// 	}
// 	return "", errors.New("Invalid string value")
// }
//
// func (v Value) Int8() (int8, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseInt(string(v.Data), 10, 8)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = int8(value)
// 	}
// 	return v.V.(int8), nil
// }
//
// func (v Value) Int16() (int16, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseInt(string(v.Data), 10, 16)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = int16(value)
// 	}
// 	return v.V.(int16), nil
// }
//
// func (v Value) Int32() (int32, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseInt(string(v.Data), 10, 32)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = int32(value)
// 	}
// 	return v.V.(int32), nil
// }
//
// func (v Value) Int64() (int64, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseInt(string(v.Data), 10, 64)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = value
// 	}
// 	return v.V.(int64), nil
// }
//
// func (v Value) UInt8() (uint8, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseUint(string(v.Data), 10, 8)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = uint8(value)
// 	}
// 	return v.V.(uint8), nil
// }
//
// func (v Value) UInt16() (uint16, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseUint(string(v.Data), 10, 16)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = uint16(value)
// 	}
// 	return v.V.(uint16), nil
// }
//
// func (v Value) UInt32() (uint32, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseUint(string(v.Data), 10, 32)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = uint32(value)
// 	}
// 	return v.V.(uint32), nil
// }
//
// func (v Value) UInt64() (uint64, error) {
// 	if v.V == nil {
// 		value, err := strconv.ParseUint(string(v.Data), 10, 64)
// 		if err != nil {
// 			return 0, err
// 		}
// 		v.V = value
// 	}
// 	return v.V.(uint64), nil
// }
//
// func (v Value) Struct() (*mdd.Containers, error) {
// 	if v.V == nil {
// 		containers, err := Decode([]byte(v.Data))
// 		if err != nil {
// 			return nil, err
// 		}
// 		v.V = containers
// 	}
// 	return v.V.(*mdd.Containers), nil
// }
//
// func (v Value) Decimal() (*big.Float, error) {
// 	if v.V == nil {
// 		f, ok := new(big.Float).SetString(string(v.Data))
// 		if !ok {
// 			return nil, errors.New("Invalid decimal value")
// 		}
// 		v.V = f
// 	}
// 	return v.V.(*big.Float), nil
// }
//
// func (v Value) DateTime() (*time.Time, error) {
// 	if v.V == nil {
// 		layout := "2006-01-02T15:04:05Z"
// 		dt, err := time.Parse(layout, string(v.Data))
// 		if err != nil {
// 			return nil, err
// 		}
// 		v.V = &dt
// 	}
// 	return v.V.(*time.Time), nil
// }
