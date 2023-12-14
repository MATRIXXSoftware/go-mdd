package mdd

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd/codec"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

type Field struct {
	CodecType   codec.Type
	Data        []byte
	Type        field.Type
	Value       interface{}
	IsMulti     bool
	IsContainer bool
}

func (f *Field) String() string {
	return string(f.Data)
}

func (f *Field) IntValue() int {
	if f.Value != nil {
		return f.Value.(int)
	}

	f.Value, _ = strconv.Atoi(f.String())
	return f.Value.(int)
}

func (f *Field) StringValue() string {
	return f.String()
}
