package mdd

import "github.com/matrixxsoftware/go-mdd/mdd/field"

// "math/big"
// "time"

type Value interface {
	Type() field.Type
	Serialize() []byte
	// String() (string, error)
	// Bool() (bool, error)
	// Int8() (int8, error)
	// Int16() (int16, error)
	// Int32() (int32, error)
	// Int64() (int64, error)
	// // Int128() (int128, error)
	// UInt8() (uint8, error)
	// UInt16() (uint16, error)
	// UInt32() (uint32, error)
	// UInt64() (uint64, error)
	// // UInt128() (uint128, error)
	// Struct() (*Containers, error)
	// Decimal() (*big.Float, error)
	// // Date() (*time.Time, error)
	// // Time() (*time.Time, error)
	// DateTime() (*time.Time, error)
}
