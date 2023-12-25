package mdd

type Codec interface {
	Decode([]byte) (*Containers, error)
	Encode(*Containers) ([]byte, error)

	DecodeField(*Field) (Value, error)
	EncodeField(*Field) ([]byte, error)
}
