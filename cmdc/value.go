package cmdc

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

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
	for i := range fields {
		field := fields[i]
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
		return nil, errors.New("invalid list value, first character must be '{'")
	}
	if b[len(b)-1] != '}' {
		return nil, errors.New("invalid list value, last character must be '}'")
	}
	var list [][]byte
	mark := 1
	roundMark := 0

	square := 0
	angle := 0
	round := 0
	curly := 1

	for idx := 1; idx < len(b); idx++ {
		c := b[idx]

		if round != 0 {
			if c == ')' {
				round--
			} else if roundMark == 0 {
				return nil, errors.New("invalid cMDC list, mismatch string length")
			} else if c == ':' {
				temp := b[roundMark+1 : idx]
				len, err := bytesToInt(temp)
				if err != nil {
					return nil, errors.New("invalid string length")
				}
				// reset round mark
				roundMark = 0
				// skip the string field
				idx += len
			} else if c < '0' || c > '9' {
				return nil, errors.New("invalid character '" + string(c) + "', numeric expected for string length")
			}
			continue
		}

		switch c {
		case '(':
			roundMark = idx
			round++
		case '[':
			square++
		case ']':
			square--
		case '<':
			angle++
		case '>':
			angle--
		case '{':
			curly++
		case '}':
			curly--
		case ',':
			if square == 0 && angle == 0 && curly == 1 {
				fieldData := b[mark:idx]
				mark = idx + 1
				list = append(list, fieldData)
			}
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
	for _, c := range []byte(v) {
		if c == '\\' {
			data = append(data, '\\', '\\')
		} else {
			data = append(data, c)
		}
	}
	data = append(data, ')')
	return data, nil
}

func fromDigit(ch byte) (byte, error) {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0', nil
	case ch >= 'A' && ch <= 'F':
		return ch - 'A' + 10, nil
	case ch >= 'a' && ch <= 'f':
		return ch - 'a' + 10, nil
	}
	return ch, errors.New("Invalid OctetString digit '" + string(ch) + "'. Valid digits are '0'-'9', 'A'-'F'")
}

func octetStringToByte(octet []byte) (byte, error) {
	b1, err := fromDigit(octet[1])
	if err != nil {
		return 0, err
	}
	b2, err := fromDigit(octet[2])
	if err != nil {
		return 0, err
	}
	return (b1 << 4) | b2, nil
}

func decodeStringWithEscapeChar(b []byte) (string, error) {
	var result []byte
	for i := 0; i < len(b); i++ {
		c := b[i]
		if c == '\\' && len(b)-i >= 3 {
			octet := b[i : i+3]
			b, err := octetStringToByte(octet)
			// Replace the char with the octet converted byte value if no error
			if err == nil {
				c = b
				i += 2
			}
		}
		result = append(result, c)
	}
	return string(result), nil
}

func decodeStringValue(b []byte) (string, error) {
	if len(b) == 0 {
		return "", nil
	}
	if b[0] != '(' {
		return "", errors.New("invalid string value")
	}
	if b[len(b)-1] != ')' {
		return "", errors.New("invalid string value")
	}
	for idx := 1; idx < len(b); idx++ {
		c := b[idx]
		if c == ':' {
			temp := b[1:idx]
			_, err := bytesToInt(temp)
			if err != nil {
				return "", errors.New("invalid string length")
			}
			result, err := decodeStringWithEscapeChar(b[idx+1 : len(b)-1])
			if err != nil {
				return "", err
			}
			return result, nil
		}
	}
	return "", errors.New("invalid string value")
}

func encodeFieldKeyValue(v string) ([]byte, error) {
	return []byte(v), nil
}

func decodeFieldKeyValue(b []byte) (string, error) {
	return string(b), nil
}

func toDigit(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'A' + (b - 10)
}

func encodeBlobValue(v string) ([]byte, error) {
	data := make([]byte, 0, len(v)*3)
	length := 0

	for _, c := range []byte(v) {
		if c == '\\' {
			data = append(data, '\\', '\\')
			length++
		} else if c > 0x1F && c < 0x7F {
			data = append(data, c)
			length++
		} else {
			data = append(data, '\\', toDigit((c&0xF0)>>4), toDigit(c&0x0F))
			length++
		}
	}

	result := make([]byte, 0, len(data)+6)
	result = append(result, '(')
	result = append(result, []byte(fmt.Sprintf("%d", length))...)
	result = append(result, ':')
	result = append(result, data...)
	result = append(result, ')')

	return result, nil
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

func encodeInt128Value(v *big.Int) ([]byte, error) {
	return []byte(v.String()), nil
}

func decodeInt128Value(b []byte) (*big.Int, error) {
	v, ok := new(big.Int).SetString(string(b), 10)
	if !ok {
		return nil, errors.New("invalid int128 value")
	}
	return v, nil
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

func encodeUInt128Value(v *big.Int) ([]byte, error) {
	return []byte(v.String()), nil
}

func decodeUInt128Value(b []byte) (*big.Int, error) {
	v, ok := new(big.Int).SetString(string(b), 10)
	if !ok {
		return nil, errors.New("invalid uint128 value")
	}
	return v, nil
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
		return nil, errors.New("invalid decimal value")
	}
	return f, nil
}

func encodeDateTimeValue(v *time.Time) ([]byte, error) {
	return []byte(v.Format(time.RFC3339)), nil
}

func decodeDateTimeValue(b []byte) (*time.Time, error) {
	layout := "2006-01-02T15:04:05Z"
	dt, err := time.Parse(layout, string(b))
	if err != nil {
		return nil, err
	}
	return &dt, nil
}

func encodeDateValue(v *time.Time) ([]byte, error) {
	return []byte(v.Format("2006-01-02")), nil
}

func decodeDateValue(b []byte) (*time.Time, error) {
	layout := "2006-01-02"
	dt, err := time.Parse(layout, string(b))
	if err != nil {
		return nil, err
	}
	return &dt, nil
}

func encodeTimeValue(v *time.Time) ([]byte, error) {
	return []byte(v.Format("15:04:05")), nil
}

func decodeTimeValue(b []byte) (*time.Time, error) {
	layout := "15:04:05"
	dt, err := time.Parse(layout, string(b))
	if err != nil {
		return nil, err
	}
	return &dt, nil
}
