package cmdc

import (
	"errors"

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
	default:
		return nil, errors.New("Unsupported field type")
	}
}

func (cmdc *Cmdc) EncodeField(f *mdd.Field) ([]byte, error) {
	// If the f has data, use it
	if len(f.Data) > 0 || f.Type == field.Unknown {
		return f.Data, nil
	}
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
	}

	return f.Data, nil
}
