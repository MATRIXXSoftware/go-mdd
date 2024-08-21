package cmdc

import (
	"fmt"
	"math/big"
	"time"

	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

type Cmdc struct {
	dict *dictionary.Dictionary
}

func NewCodecWithDict(dict *dictionary.Dictionary) *Cmdc {
	return &Cmdc{
		dict: dict,
	}
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
		case field.Int128:
			return decodeInt128Value(f.Data)
		case field.UInt8:
			return decodeUInt8Value(f.Data)
		case field.UInt16:
			return decodeUInt16Value(f.Data)
		case field.UInt32:
			return decodeUInt32Value(f.Data)
		case field.UInt64:
			return decodeUInt64Value(f.Data)
		case field.UInt128:
			return decodeUInt128Value(f.Data)
		case field.Bool:
			return decodeBoolValue(f.Data)
		case field.Struct:
			return decodeStructValue(f.Codec, f.Data)
		case field.Decimal:
			return decodeDecimalValue(f.Data)
		case field.Date:
			return decodeDateValue(f.Data)
		case field.Time:
			return decodeTimeValue(f.Data)
		case field.DateTime:
			return decodeDateTimeValue(f.Data)
		case field.Blob:
			return decodeStringValue(f.Data)
		case field.BufferID:
			return decodeStringValue(f.Data)
		case field.FieldKey:
			return decodeFieldKeyValue(f.Data)
		case field.PhoneNo:
			return field.NewPhoneNo(f.Data)
		case field.ObjectID:
			return field.NewObjectID(f.Data)
		default:
			return nil, fmt.Errorf("unsupported field type '%v'", f.Type)
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
		case field.Int128:
			return decodeListValue(f.Data, decodeInt128Value)
		case field.UInt8:
			return decodeListValue(f.Data, decodeUInt8Value)
		case field.UInt16:
			return decodeListValue(f.Data, decodeUInt16Value)
		case field.UInt32:
			return decodeListValue(f.Data, decodeUInt32Value)
		case field.UInt64:
			return decodeListValue(f.Data, decodeUInt64Value)
		case field.UInt128:
			return decodeListValue(f.Data, decodeUInt128Value)
		case field.Bool:
			return decodeListValue(f.Data, decodeBoolValue)
		case field.Struct:
			return decodeListValue(f.Data, func(b []byte) (*mdd.Containers, error) {
				return decodeStructValue(f.Codec, b)
			})
		case field.Decimal:
			return decodeListValue(f.Data, decodeDecimalValue)
		case field.Date:
			return decodeListValue(f.Data, decodeDateValue)
		case field.Time:
			return decodeListValue(f.Data, decodeTimeValue)
		case field.DateTime:
			return decodeListValue(f.Data, decodeDateTimeValue)
		case field.Blob:
			return decodeListValue(f.Data, decodeStringValue)
		case field.BufferID:
			return decodeListValue(f.Data, decodeStringValue)
		case field.FieldKey:
			return decodeListValue(f.Data, decodeFieldKeyValue)
		case field.PhoneNo:
			return decodeListValue(f.Data, field.NewPhoneNo)
		case field.ObjectID:
			return decodeListValue(f.Data, field.NewObjectID)
		default:
			return nil, fmt.Errorf("unsupported field type '%v'", f.Type)
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
		case field.Int128:
			return encodeInt128Value(f.Value.(*big.Int))
		case field.UInt8:
			return encodeUInt8Value(f.Value.(uint8))
		case field.UInt16:
			return encodeUInt16Value(f.Value.(uint16))
		case field.UInt32:
			return encodeUInt32Value(f.Value.(uint32))
		case field.UInt64:
			return encodeUInt64Value(f.Value.(uint64))
		case field.UInt128:
			return encodeUInt128Value(f.Value.(*big.Int))
		case field.Bool:
			return encodeBoolValue(f.Value.(bool))
		case field.Struct:
			return encodeStructValue(f.Codec, f.Value.(*mdd.Containers))
		case field.Decimal:
			return encodeDecimalValue(f.Value.(*big.Float))
		case field.Date:
			return encodeDateValue(f.Value.(*time.Time))
		case field.Time:
			return encodeTimeValue(f.Value.(*time.Time))
		case field.DateTime:
			return encodeDateTimeValue(f.Value.(*time.Time))
		case field.Blob:
			return encodeBlobValue(f.Value.(string))
		case field.BufferID:
			return encodeStringValue(f.Value.(string))
		case field.FieldKey:
			return encodeFieldKeyValue(f.Value.(string))
		case field.PhoneNo:
			return f.Value.(field.MtxPhoneNo).Bytes()
		case field.ObjectID:
			return f.Value.(field.MtxObjectID).Bytes()
		default:
			return nil, fmt.Errorf("unsupported field type '%v'", f.Type)
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
		case field.Int128:
			return encodeListValue(f.Value.([]*big.Int), encodeInt128Value)
		case field.UInt8:
			return encodeListValue(f.Value.([]uint8), encodeUInt8Value)
		case field.UInt16:
			return encodeListValue(f.Value.([]uint16), encodeUInt16Value)
		case field.UInt32:
			return encodeListValue(f.Value.([]uint32), encodeUInt32Value)
		case field.UInt64:
			return encodeListValue(f.Value.([]uint64), encodeUInt64Value)
		case field.UInt128:
			return encodeListValue(f.Value.([]*big.Int), encodeUInt128Value)
		case field.Bool:
			return encodeListValue(f.Value.([]bool), encodeBoolValue)
		case field.Struct:
			return encodeListValue(f.Value.([]*mdd.Containers), func(v *mdd.Containers) ([]byte, error) {
				return encodeStructValue(f.Codec, v)
			})
		case field.Decimal:
			return encodeListValue(f.Value.([]*big.Float), encodeDecimalValue)
		case field.Date:
			return encodeListValue(f.Value.([]*time.Time), encodeDateValue)
		case field.Time:
			return encodeListValue(f.Value.([]*time.Time), encodeTimeValue)
		case field.DateTime:
			return encodeListValue(f.Value.([]*time.Time), encodeDateTimeValue)
		case field.Blob:
			return encodeListValue(f.Value.([]string), encodeBlobValue)
		case field.BufferID:
			return encodeListValue(f.Value.([]string), encodeStringValue)
		case field.FieldKey:
			return encodeListValue(f.Value.([]string), encodeFieldKeyValue)
		case field.PhoneNo:
			return encodeListValue(f.Value.([]field.MtxPhoneNo), func(v field.MtxPhoneNo) ([]byte, error) {
				return v.Bytes()
			})
		case field.ObjectID:
			return encodeListValue(f.Value.([]field.MtxObjectID), func(v field.MtxObjectID) ([]byte, error) {
				return v.Bytes()
			})
		default:
			return nil, fmt.Errorf("unsupported field type '%v'", f.Type)
		}
	}
}
