package mdd

import (
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

type Field struct {
	Data        []byte
	Type        field.Type
	Value       Value
	IsMulti     bool
	IsContainer bool
}

func (f *Field) String() string {
	return string(f.Data)
}

type Value interface {
	String() (string, error)
	// Bool() (bool, error)
	// Int8() (int8, error)
	// Int16() (int16, error)
	Int32() (int32, error)
	// Int64() (int64, error)
	// Int128() (int128, error)
	// UInt8() (uint8, error)
	// UInt16() (uint16, error)
	// UInt32() (uint32, error)
	// UInt64() (uint64, error)
	// UInt128() (uint128, error)
	Float32() (float32, error)
	// Float64() (float64, error)
	Struct() (*Containers, error)
}
