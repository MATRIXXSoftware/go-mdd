package cmdc

import (
	"errors"
	"strconv"

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

func (c *Cmdc) DecodeField(f *mdd.Field) (mdd.Value, error) {
	switch f.Type {
	case field.String:
		return DecodeStringValue(f.Data)
	case field.Int32:
		return DecodeInt32Value(f.Data)
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
	case field.Int32:
		v := f.Value.(Int32Value)
		// v, err := f.Value.Int32()
		// if err != nil {
		// 	return nil, err
		// }
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case field.String:
		v := f.Value.(StringValue)
		// v, err := f.Value.String()
		// if err != nil {
		// 	return nil, err
		// }
		data := make([]byte, 0, len(v)+6)
		data = append(data, '(')
		data = append(data, []byte(strconv.Itoa(len(v)))...)
		data = append(data, ':')
		data = append(data, []byte(v)...)
		data = append(data, ')')
		return data, nil
	}

	return f.Data, nil
}
