package mdd

import (
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
