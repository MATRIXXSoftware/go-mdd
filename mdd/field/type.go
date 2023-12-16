package field

type Type int

const (
	Unknown Type = iota
	String
	Bool
	Int8
	Int16
	Int32
	Int64
	Int128
	UInt8
	UInt16
	UInt32
	UInt64
	UInt128
	Struct
)

func (t Type) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case String:
		return "String"
	case Bool:
		return "Bool"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Int128:
		return "Int128"
	case UInt8:
		return "UInt8"
	case UInt16:
		return "UInt16"
	case UInt32:
		return "UInt32"
	case UInt64:
		return "UInt64"
	case UInt128:
		return "UInt128"
	case Struct:
		return "Struct"
	default:
		return "Undefined"
	}
}
