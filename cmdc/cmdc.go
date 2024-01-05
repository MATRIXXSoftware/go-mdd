package cmdc

import (
	"errors"
	"math/big"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

type Cmdc struct {
}

func NewCodec() *Cmdc {
	return &Cmdc{}
}

func (c *Cmdc) Decode(data []byte) (*mdd.Containers, error) {
	return c.decodeContainers(data)
}

func (c *Cmdc) Encode(containers *mdd.Containers) ([]byte, error) {
	return c.encodeContainers(containers)
}

func (c *Cmdc) DecodeField(f *mdd.Field) (interface{}, error) {
	if !f.IsMulti {
		switch f.Type {
		case field.String:
			return decodeStringValue(f.Data)
		case field.Int8:
			return decodeInt8Value(f.Data)
		case field.Int16:
			return decodeInt16Value(f.Data)
		case field.Int32:
			return decodeInt32Value(f.Data)
		case field.Int64:
			return decodeInt64Value(f.Data)
		case field.UInt8:
			return decodeUInt8Value(f.Data)
		case field.UInt16:
			return decodeUInt16Value(f.Data)
		case field.UInt32:
			return decodeUInt32Value(f.Data)
		case field.UInt64:
			return decodeUInt64Value(f.Data)
		case field.Bool:
			return decodeBoolValue(f.Data)
		case field.Struct:
			return decodeStructValue(f.Codec, f.Data)
		case field.Decimal:
			return decodeDecimalValue(f.Data)
		default:
			return nil, errors.New("Unsupported field type")
		}
	} else {
		switch f.Type {
		case field.String:
			return decodeListValue(f.Data, decodeStringValue)
		case field.Int8:
			return decodeListValue(f.Data, decodeInt8Value)
		case field.Int16:
			return decodeListValue(f.Data, decodeInt16Value)
		case field.Int32:
			return decodeListValue(f.Data, decodeInt32Value)
		case field.Int64:
			return decodeListValue(f.Data, decodeInt64Value)
		case field.UInt8:
			return decodeListValue(f.Data, decodeUInt8Value)
		case field.UInt16:
			return decodeListValue(f.Data, decodeUInt16Value)
		case field.UInt32:
			return decodeListValue(f.Data, decodeUInt32Value)
		case field.UInt64:
			return decodeListValue(f.Data, decodeUInt64Value)
		case field.Bool:
			return decodeListValue(f.Data, decodeBoolValue)
		case field.Struct:
			return decodeListValue(f.Data, func(b []byte) (*mdd.Containers, error) {
				return decodeStructValue(f.Codec, b)
			})
		case field.Decimal:
			return decodeListValue(f.Data, decodeDecimalValue)
		default:
			return nil, errors.New("Unsupported field type")
		}
	}
}

func (cmdc *Cmdc) EncodeField(f *mdd.Field) ([]byte, error) {
	// If the f is null, return empty data
	if f.IsNull {
		return []byte{}, nil
	}
	// If the f has data, use it
	if len(f.Data) > 0 || f.Type == field.Unknown {
		return f.Data, nil
	}

	if !f.IsMulti {
		switch f.Type {
		case field.String:
			return encodeStringValue(f.Value.(string))
		case field.Int8:
			return encodeInt8Value(f.Value.(int8))
		case field.Int16:
			return encodeInt16Value(f.Value.(int16))
		case field.Int32:
			return encodeInt32Value(f.Value.(int32))
		case field.Int64:
			return encodeInt64Value(f.Value.(int64))
		case field.UInt8:
			return encodeUInt8Value(f.Value.(uint8))
		case field.UInt16:
			return encodeUInt16Value(f.Value.(uint16))
		case field.UInt32:
			return encodeUInt32Value(f.Value.(uint32))
		case field.UInt64:
			return encodeUInt64Value(f.Value.(uint64))
		case field.Bool:
			return encodeBoolValue(f.Value.(bool))
		case field.Struct:
			return encodeStructValue(f.Codec, f.Value.(*mdd.Containers))
		case field.Decimal:
			return encodeDecimalValue(f.Value.(*big.Float))
		default:
			return nil, errors.New("Unsupported field type")
		}
	} else {
		switch f.Type {
		case field.String:
			return encodeListValue(f.Value.([]string), encodeStringValue)
		case field.Int8:
			return encodeListValue(f.Value.([]int8), encodeInt8Value)
		case field.Int16:
			return encodeListValue(f.Value.([]int16), encodeInt16Value)
		case field.Int32:
			return encodeListValue(f.Value.([]int32), encodeInt32Value)
		case field.Int64:
			return encodeListValue(f.Value.([]int64), encodeInt64Value)
		case field.UInt8:
			return encodeListValue(f.Value.([]uint8), encodeUInt8Value)
		case field.UInt16:
			return encodeListValue(f.Value.([]uint16), encodeUInt16Value)
		case field.UInt32:
			return encodeListValue(f.Value.([]uint32), encodeUInt32Value)
		case field.UInt64:
			return encodeListValue(f.Value.([]uint64), encodeUInt64Value)
		case field.Bool:
			return encodeListValue(f.Value.([]bool), encodeBoolValue)
		case field.Struct:
			return encodeListValue(f.Value.([]*mdd.Containers), func(v *mdd.Containers) ([]byte, error) {
				return encodeStructValue(f.Codec, v)
			})
		case field.Decimal:
			return encodeListValue(f.Value.([]*big.Float), encodeDecimalValue)
		default:
			// TODO: Add support for other types
			return nil, errors.New("Unsupported field type")
		}
	}
}
