package cmdc

import (
	"errors"
	"strconv"
)

func decodeStringValue(b []byte) (string, error) {
	if len(b) == 0 {
		return string(""), nil
	}
	if b[0] != '(' {
		return string(""), errors.New("Invalid string value")
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
			return string(str), nil
		}
	}
	return string(""), errors.New("Invalid string value")
}

func decodeInt32Value(b []byte) (int32, error) {
	v, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return int32(0), err
	}
	return int32(v), nil
}

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
