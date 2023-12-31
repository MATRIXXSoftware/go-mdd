package cmdc

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func encodeListValue[T any](list []T, f func(T) ([]byte, error)) ([]byte, error) {
	// todo estimate size to allocate
	data := make([]byte, 0, len(list)*8+2)
	data = append(data, '{')

	if len(list) > 0 {
		// first element
		b, err := f(list[0])
		if err != nil {
			return nil, err
		}
		data = append(data, b...)

		// remaining elements
		for i := 1; i < len(list); i++ {
			data = append(data, ',')
			b, err := f(list[i])
			if err != nil {
				return nil, err
			}
			data = append(data, b...)
		}
	}

	data = append(data, '}')
	return data, nil
}

func decodeListValue[T any](b []byte, f func([]byte) (T, error)) ([]T, error) {
	fields, err := decodeList(b)
	if err != nil {
		return nil, err
	}
	var list []T
	for _, field := range fields {
		v, err := f(field)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func decodeList(b []byte) ([][]byte, error) {
	if len(b) == 0 {
		return nil, nil
	}
	if b[0] != '{' {
		return nil, errors.New("Invalid list value, first character must be '{'")
	}
	if b[len(b)-1] != '}' {
		return nil, errors.New("Invalid list value, last character must be '}'")
	}
	var list [][]byte
	mark := 1
	for idx := 1; idx < len(b); idx++ {
		c := b[idx]
		if c == ',' {
			fieldData := b[mark:idx]
			mark = idx + 1
			list = append(list, fieldData)
		}
	}
	fieldData := b[mark : len(b)-1]
	list = append(list, fieldData)
	return list, nil
}

func encodeBoolValue(v bool) ([]byte, error) {
	if v == true {
		return []byte("1"), nil
	} else {
		return []byte("0"), nil
	}

}

func decodeBoolValue(b []byte) (bool, error) {
	v, err := strconv.ParseBool(string(b))
	if err != nil {
		return false, err
	}
	return v, nil
}

func encodeStringValue(v string) ([]byte, error) {
	data := make([]byte, 0, len(v)+6)
	data = append(data, '(')
	data = append(data, []byte(strconv.Itoa(len(v)))...)
	data = append(data, ':')
	data = append(data, []byte(v)...)
	data = append(data, ')')
	return data, nil
}

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

func encodeStringListValue(v []string) ([]byte, error) {
	return encodeListValue(v, encodeStringValue)
}

func decodeStringListValue(b []byte) ([]string, error) {
	return decodeListValue(b, decodeStringValue)
}

func encodeInt8Value(v int8) ([]byte, error) {
	return []byte(strconv.FormatInt(int64(v), 10)), nil
}

func decodeInt8Value(b []byte) (int8, error) {
	v, err := strconv.ParseInt(string(b), 10, 8)
	if err != nil {
		return int8(0), err
	}
	return int8(v), nil
}

func encodeInt16Value(v int16) ([]byte, error) {
	return []byte(strconv.FormatInt(int64(v), 10)), nil
}

func decodeInt16Value(b []byte) (int16, error) {
	v, err := strconv.ParseInt(string(b), 10, 16)
	if err != nil {
		return int16(0), err
	}
	return int16(v), nil
}

func encodeInt32Value(v int32) ([]byte, error) {
	return []byte(strconv.FormatInt(int64(v), 10)), nil
}

func decodeInt32Value(b []byte) (int32, error) {
	v, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return int32(0), err
	}
	return int32(v), nil
}

func encodeInt32ListValue(v []int32) ([]byte, error) {
	return encodeListValue(v, encodeInt32Value)
}

func decodeInt32ListValue(b []byte) ([]int32, error) {
	return decodeListValue(b, decodeInt32Value)
}

func encodeInt64Value(v int64) ([]byte, error) {
	return []byte(strconv.FormatInt(v, 10)), nil
}

func decodeInt64Value(b []byte) (int64, error) {
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return int64(0), err
	}
	return int64(v), nil
}

func encodeUInt8Value(v uint8) ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(v), 10)), nil
}

func decodeUInt8Value(b []byte) (uint8, error) {
	v, err := strconv.ParseUint(string(b), 10, 8)
	if err != nil {
		return uint8(0), err
	}
	return uint8(v), nil
}

func encodeUInt16Value(v uint16) ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(v), 10)), nil
}

func decodeUInt16Value(b []byte) (uint16, error) {
	v, err := strconv.ParseUint(string(b), 10, 16)
	if err != nil {
		return uint16(0), err
	}
	return uint16(v), nil
}

func encodeUInt32Value(v uint32) ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(v), 10)), nil
}

func decodeUInt32Value(b []byte) (uint32, error) {
	v, err := strconv.ParseUint(string(b), 10, 32)
	if err != nil {
		return uint32(0), err
	}
	return uint32(v), nil
}

func encodeUInt64Value(v uint64) ([]byte, error) {
	return []byte(strconv.FormatUint(v, 10)), nil
}

func decodeUInt64Value(b []byte) (uint64, error) {
	v, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return uint64(0), err
	}
	return uint64(v), nil
}

func encodeStructValue(codec mdd.Codec, v *mdd.Containers) ([]byte, error) {
	return codec.Encode(v)
}

func decodeStructValue(codec mdd.Codec, b []byte) (*mdd.Containers, error) {
	return codec.Decode(b)
}

func encodeDecimalValue(v *big.Float) ([]byte, error) {
	return []byte(v.String()), nil
}

func decodeDecimalValue(b []byte) (*big.Float, error) {
	f, ok := new(big.Float).SetString(string(b))
	if !ok {
		return nil, errors.New("Invalid decimal value")
	}
	return f, nil
}

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
